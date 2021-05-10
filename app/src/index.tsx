import React, {FC} from "react";
import ReactDOM from "react-dom";
import "./styles.scss";
import App from "./Containers/App";
import {Router} from "react-router-dom";
import history from "./history";
import {WebhookProvider} from "./Contexts/WebHookContext";

const RootApp: FC = () => {
  return (
    <WebhookProvider>
      <Router history={history}>
        <App/>
      </Router>
    </WebhookProvider>
  )
}

ReactDOM.render(<RootApp/>, document.getElementById("root"))