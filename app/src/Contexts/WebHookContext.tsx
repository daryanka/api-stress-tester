import React, {createContext, FC, useContext, useState} from "react";

interface contextState {
  data: any,
  setData: any,
  con: WebSocket | undefined,
  setCon: React.Dispatch<React.SetStateAction<WebSocket | undefined>>
}

const WebhookContext = createContext<any>(null)

export const WebhookProvider: FC = (props) => {
  const [data, setData] = useState({
    test: []
  })
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
    const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6NywiZXhwIjoxNjMxMDMyNDg4LCJpYXQiOjE2MjEwMzI0ODgsImlzcyI6ImpvdXJuYWxfYXBpIn0.2hH2M1eCRDiDbGYOllA98f7p4hr0EfV4s3BYZh6dUy4"
    const socket = new WebSocket("ws://localhost:8081/v1/ws", token)
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
    setData((prev: any) => ({
      test: [...prev.test, msg]
    }))
  }

  const send = (data: any) => {
    if (con) {
      con.send(JSON.stringify(data))
    } else {
      console.log("Unable to send message no socket connection open")
    }
  }


  return {
    initialiseConnection,
    data,
    send
  }
}