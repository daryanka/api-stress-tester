import React, {FC} from "react";
import ReactDOM from "react-dom";
import "./styles.scss";
import App from "./Containers/App";
import {Router} from "react-router-dom";
import history from "./history";

const RootApp: FC = () => {
  return (
    <Router history={history}>
      <App/>
    </Router>
  )
}

ReactDOM.render(<RootApp/>, document.getElementById("root"))