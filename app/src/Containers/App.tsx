import React, {FC, useState} from "react";
import {useWebhook} from "../Contexts/WebHookContext";

const App: FC = () => {
  console.log(process.env)
  const {initialiseConnection, send, data} = useWebhook()

  const msg = () => {
    send({
      "test": true,
      "message" : "hello world"
    })
  }

  const handleConnect = () => {
    initialiseConnection()
  }

  return (
    <div>
      <h1>Typescript Template</h1>
      <button onClick={handleConnect}>Connect</button>
      <button onClick={msg}>Send</button>
      {data.test && data.test.length}
    </div>
  )
}

export default App