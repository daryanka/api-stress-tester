import React, {FC} from "react";
import {Route, Redirect, Switch} from "react-router-dom";
import Login from "./Login";
import Nav from "../Components/Nav";
import PrivateRoute from "../Components/PrivateRoute";
import Register from "./Register";
import RequestPage from "./requests/Requests";

const App: FC = () => {
  return (
    <div className={"app-wrapper"}>
      <Nav />
        <Switch>
          <Route path={"/"} exact component={Temp}/>
          <Route path={"/login"} exact component={Login}/>
          <Route path={"/register"} exact component={Register}/>
          <PrivateRoute path={"/dashboard"} exact />
          <PrivateRoute path={"/domains"} exact />
          <PrivateRoute path={"/stress-test"} component={RequestPage} />
          <Redirect to={"/"} />
        </Switch>
    </div>
  )
}

const Temp: FC = () => {
  return (
    <div>
      <h1>Home</h1>
    </div>
  )
}

export default App