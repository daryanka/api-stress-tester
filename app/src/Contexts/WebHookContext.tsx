import React, {createContext, FC, useContext, useState} from "react";

interface contextState {
  data: any,
  setData: any,
  con: WebSocket | undefined,
  setCon: React.Dispatch<React.SetStateAction<WebSocket | undefined>>
}

const WebhookContext = createContext<any>(null)

export const WebhookProvider: FC = (props) => {
  const [data, setData] = useState({})
  const [con, setCon] = useState<WebSocket>()

  return (
    <WebhookContext.Provider value={{data, setData, con, setCon}}>
      {props.children}
    </WebhookContext.Provider>
  )
}

export const useWebhook = () => {
  const {data, setData, con, setCon} = useContext(WebhookContext) as contextState

  const initialiseConnection = () => {
    // Send
    const socket = new WebSocket("ws://localhost:8081/v1/ws", "testabc")
    socket.onopen = (ev) => {
      console.log(ev)
      setCon(socket)
    }

    socket.onclose = (ev) => {
      setCon(undefined)
      console.log("connection closed", ev)
    }

    socket.onmessage = (msg) => {
      handleMessage(msg)
    }
  }

  const handleMessage = (msg :MessageEvent) => {

  }

  const send = (data: any) => {
    if (con) {
      con.send(data)
    } else {
      console.log("Unable to send message no socket connection open")
    }
  }


  return {
    initialiseConnection,
    data
  }
}