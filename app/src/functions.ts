import axios, {AxiosRequestConfig, AxiosResponse} from "axios";
import cookie from "js-cookie";
import _ from "lodash";
import {useHistory} from "react-router-dom";

type Methods = "GET" | "PUT" | "PATCH" | "POST" | "DELETE";

export interface ErrorType {
  error: boolean
  message: string
  type?: string
  status_code: number
}

const useAPI = () => {
  const history = useHistory();

  const send = async <T = any>(method: Methods, url: string, data?: any, additionalConfig?: AxiosRequestConfig) => {
    const headers: {
      Authorization?: string
    } = {};

    if (cookie.get("token")) {
      headers.Authorization = `Bearer ${cookie.get("token")}`
    }

    let baseURL = ""

    switch (process.env.NODE_ENV) {
      case "dev":
        baseURL = `http://localhost:8080${url}`;
        break;
      case "prod":
        // TODO add prod when ready
        baseURL = `http://localhost:8080${url}`;
        break;
      default:
        baseURL = `http://localhost:8080${url}`;
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

  const get = <T = any>(url: string, config?: AxiosRequestConfig) => {
    return send<T>("GET", url, undefined, config)
  }

  const post = <T = any>(url: string, data: any, config?: AxiosRequestConfig) => {
    return send<T>("POST", url, data, config)
  }

  const patch = <T = any>(url: string, data: object, config?: AxiosRequestConfig) => {
    return send<T>("PATCH", url, data, config)
  }

  const put = <T = any>(url: string, data: object, config?: AxiosRequestConfig) => {
    return send<T>("PUT", url, data, config)
  }

  const deleteFunc = <T = any>(url: string, data?: object, config?: AxiosRequestConfig) => {
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

  return {
    send,
    get,
    post,
    patch,
    put,
    delete: deleteFunc,
    error
  }
}

export default useAPI