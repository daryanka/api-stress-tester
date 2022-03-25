import axios, {AxiosRequestConfig, AxiosResponse} from "axios";
import cookie from "js-cookie";
import _ from "lodash";
import {useHistory} from "react-router-dom";
import {FormikHelpers} from "formik";
import {useQueryClient} from "react-query";
import {useAuthenticated} from "./Contexts/AuthenticationContext";
import {useWebhook} from "./Contexts/WebHookContext";

type Methods = "GET" | "PUT" | "PATCH" | "POST" | "DELETE";

export interface ErrorType {
  error: string
  type?: string
  status_code: number
}

const useAPI = () => {
  const history = useHistory();
  const queryClient = useQueryClient()
  const {setIsAuthenticated} = useAuthenticated()
  const {closeConnection} = useWebhook()

  const send = async <T>(method: Methods, url: string, data?: any, additionalConfig?: AxiosRequestConfig) => {
    const headers: {
      Authorization?: string
    } = {};

    if (cookie.get("token")) {
      headers.Authorization = `Bearer ${cookie.get("token")}`
    }

    let baseURL = ""

    switch (process.env.NODE_ENV) {
      case "DEV":
        baseURL = `http://localhost:8081/v1${url}`;
        break;
      case "PROD":
        baseURL = `https://api-api-tester.daryanamin.co.uk/v1${url}`;
        break;
      default:
        baseURL = `http://localhost:8081/v1${url}`;
        break;
    }

    baseURL = `https://api-api-tester.daryanamin.co.uk/v1${url}`;

    const config = _.merge(additionalConfig ? additionalConfig : {}, {
      method: method,
      url: baseURL,
      headers: headers,
      data: data,
      validateStatus: (status: number) => {
        // return status >= 200 && status < 300
        return true
      }
    });

    return await axios.request<T>(config)
  }

  const get = <T>(url: string, config?: AxiosRequestConfig) => {
    return send<T>("GET", url, undefined, config)
  }

  const post = <T>(url: string, data: any, config?: AxiosRequestConfig) => {
    return send<T>("POST", url, data, config)
  }

  const patch = <T>(url: string, data: object, config?: AxiosRequestConfig) => {
    return send<T>("PATCH", url, data, config)
  }

  const put = <T>(url: string, data: object, config?: AxiosRequestConfig) => {
    return send<T>("PUT", url, data, config)
  }

  const deleteFunc = <T>(url: string, data?: object, config?: AxiosRequestConfig) => {
    return send<T>("DELETE", url, data, config)
  }

  const error = (p: AxiosResponse) => {
    if (p.status >= 300 || p.status < 200) {
      // Error
      if (p.status === 401) {
        // Logout
        pushTo("/logout")
      }

      return p.data as ErrorType
    }
    return null
  }

  const pushTo = (url: string) => {
    history.push(url)
  }

  const handleFormikError = (res: AxiosResponse, helpers: FormikHelpers<any>) => {
    if (res.status < 200 || res.status >= 300) {
      // Has error
      // Check if validation error
      if (res.status === 422) {
        // Add errors to obj
        const obj: { [key: string]: string } = {}
        const keys = Object.keys(res.data)
        for (let i = 0; i < keys.length; i++) {
          obj[keys[i]] = res.data[keys[i]]
        }

        console.log("status", obj)
        helpers.setStatus(obj)
        return true
      }

      // set error
      helpers.setStatus({
        DEFAULT_ERROR: res?.data?.error
      })

      helpers.setTouched({})

      return true
    }

    return false
  }

  const logout = () => {
    cookie.remove("token")
    queryClient.invalidateQueries(["me"])
    history.replace("/login")
    setIsAuthenticated(false)
    closeConnection()
  }

  return {
    send,
    get,
    post,
    patch,
    put,
    delete: deleteFunc,
    handleFormikError,
    error,
    logout
  }
}

/**
 * Seconds / Minutes string to seconds number
 */
export const SMToMinutes = (str: string) => {
  // Currently only working with minutes and seconds, can be configured to work with hours, days and weeks as well.

  // To lowercase and remove all consecutive spaces to single space
  str = str.toLowerCase().replace(/\s+/g,' ')

  // Split by spaces and remove all empty strings or just spaces
  const timesArr = str.split(" ").filter(x => {
    // If returns true keep
    return !(x === " " || x === "");
  })

  const time = {
    m: 0,
    s: 0,
    // h: 0,
    // d: 0,
    // w: 0,
  }

  // type timeTypes = "m" | "h" | "d" | "w" | "s";
  type timeTypes = "m" | "s";

  let valid: boolean = true;
  for (let i = 0; i < timesArr.length; i++) {
    // const regex = /^([0-9]+)([wdmhs])$/;
    const regex = /^([0-9]+)([ms])$/;
    if (!regex.test(timesArr[i])) {
      //Invalid
      valid = false;
      break;
    } else {
      //Valid
      const match = regex.exec(timesArr[i])
      const number: number = parseInt(match![1])
      const timeType: timeTypes = match![2] as timeTypes;

      time[timeType] = time[timeType] + number;
    }
  }

  let totalSeconds = 0;

  for (let key in time) {
    switch (key as timeTypes) {
      case "s":
        totalSeconds = totalSeconds + time.s
        break;
      case "m":
        totalSeconds = totalSeconds + (time.m * 60)
        break;
      // case "h":
      //   totalSeconds = totalSeconds + (60 * 60 * time.h)
      //   break;
      // case "d":
      //   totalSeconds = totalSeconds + (60 * 24 * time.d)
      //   break;
      // case "w":
      //   totalSeconds = totalSeconds + (60 * 24 * 7 * time.w)
      //   break;
    }
  }

  let PrettyString = []
  let temp = totalSeconds

  if (totalSeconds / 60 > 1) {
    PrettyString.push(`${Math.floor(totalSeconds / 60)} Minutes`)
    temp = temp - Math.floor(totalSeconds / 60) * 60
  }
  if (temp > 0) {
    PrettyString.push(`${temp} Seconds`)
  }

  return {
    valid,
    totalSeconds,
    PrettyString: PrettyString.join(", ")
  }
}

export const SecondsToPretty = (seconds: number) => {
  let remaining = seconds
  const PrettyString = []

  if (remaining / 60 > 1) {
    PrettyString.push(`${Math.floor(remaining / 60)} Minutes`)
    remaining = remaining - Math.floor(remaining / 60) * 60
  }
  if (remaining > 0) {
    PrettyString.push(`${remaining} Seconds`)
  }

  return PrettyString.join(", ")
}

export default useAPI