import React, {FC, useState} from "react";
import {useWebhook} from "../Contexts/WebHookContext";

const App: FC = () => {
  const {initialiseConnection} = useWebhook()

  const handleConnect = () => {
    initialiseConnection()
  }

  return (
    <div>
      <h1>Typescript Template</h1>
      <button onClick={handleConnect}>Connect</button>
    </div>
  )
}

export default App