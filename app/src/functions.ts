import axios, {AxiosRequestConfig, AxiosResponse} from "axios";
import cookie from "js-cookie";
import _ from "lodash";
import {useHistory} from "react-router-dom";
import {FormikHelpers} from "formik";
import {useQueryClient} from "react-query";
import {useAuthenticated} from "./Contexts/AuthenticationContext";

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

  const send = async <T>(method: Methods, url: string, data?: any, additionalConfig?: AxiosRequestConfig) => {
    const headers: {
      Authorization?: string
    } = {};

    if (cookie.get("token")) {
      headers.Authorization = `Bearer ${cookie.get("token")}`
    }

    let baseURL = ""

    switch (process.env.NODE_ENV) {
      case "dev":
        baseURL = `http://localhost:8081/v1${url}`;
        break;
      case "prod":
        // TODO add prod when ready
        baseURL = `http://localhost:8081/v1${url}`;
        break;
      default:
        baseURL = `http://localhost:8081/v1${url}`;
        break;
    }


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

      // Check if validatoin error
      if (res.status === 422) {
        // Add errors to obj
        const obj: { [key: string]: string } = {}
        const keys = Object.keys(res.data)
        for (let i = 0; i < keys.length; i++) {
          obj[keys[i]] = res.data[keys[i]]
        }

        helpers.setStatus(obj)
        return true
      }

      // set error
      helpers.setStatus({
        DEFAULT_ERROR: res?.data?.error
      })

      return true
    }

    return false
  }

  const logout = () => {
    cookie.remove("token")
    queryClient.invalidateQueries(["me"])
    history.replace("/login")
    setIsAuthenticated(false)
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

export default useAPI