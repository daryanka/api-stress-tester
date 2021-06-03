import React, {createContext, FC, useContext, useState} from "react";

type FullRequest =  {
  [key: string]: {
    reqs: IndividualRequest[],
    status: string,
  }
}

interface contextState {
  requests: FullRequest,
  setRequests: React.Dispatch<React.SetStateAction<FullRequest>>,
  con: WebSocket | undefined,
  setCon: React.Dispatch<React.SetStateAction<WebSocket | undefined>>
}

interface MessageType {
  type: string
  status: number
  info: any
}

interface IndividualRequest {
  request_id: number
  time_taken: number
  status_code: number
}

interface CompleteInfo {
  request_id: number
  message: string
}

const WebhookContext = createContext<any>(null)

export const WebhookProvider: FC = (props) => {
  const [requests, setRequests] = useState({})
  const [con, setCon] = useState<WebSocket>()

  return (
    <WebhookContext.Provider value={{requests, setRequests, con, setCon}}>
      {props.children}
    </WebhookContext.Provider>
  )
}

export const useWebhook = () => {
  const {setRequests, setCon} = useContext(WebhookContext) as contextState

  const initialiseConnection = () => {
    // Send
    const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6NywiZXhwIjoxNjMyNzQ2ODg4LCJpYXQiOjE2MjI3NDY4ODgsImlzcyI6ImpvdXJuYWxfYXBpIn0.9iKCH1cJRjjbrJxvHK3_10GWqBkIOoTR8dAtmVTZ_bU"
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
    const data = JSON.parse(msg.data) as MessageType

    setRequests((prev) => {
      prev = {...prev};

      switch (data.type) {
        case "INDIVIDUAL_REQUEST":
          // Add to request
          const info = data.info as IndividualRequest
          if (prev.hasOwnProperty(info.request_id)) {
            prev[info.request_id].reqs.push({...data.info})
          } else {
            prev[info.request_id] = {
              reqs: [{...data.info}],
              status: "In Progress"
            }
          }
          break
        case "REQUEST_FAILED":
        case "REQUEST_COMPLETE":
          // Signal new api request to get all data
          const failedInfo = data.info as CompleteInfo
          if (prev.hasOwnProperty(failedInfo.request_id)) {
            prev[failedInfo.request_id].status = failedInfo.message
          }
          break
        default:
          break
      }
      return prev
    })
  }


  return {
    initialiseConnection,
  }
}

export const useRequestHook = (requestID: number) => {
  const {requests, con} = useContext(WebhookContext) as contextState

  const sendCancel = () => {
    const data = {
      type: "CANCEL_REQUEST",
      requestID: requestID
    }
    if (con) {
      try {
        con.send(JSON.stringify(data))
      } catch (e) {
        console.log("Unable to send message no socket connection open")
      }
    } else {
      console.log("Unable to send message no socket connection open")
    }
  }

  return {
    fullRequest: requests[requestID],
    sendCancel,
  }
}