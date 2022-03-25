import React, {FC} from "react";
import {Route, Redirect, Switch, useHistory} from "react-router-dom";
import Login from "./Login";
import Nav from "../Components/Nav";
import PrivateRoute from "../Components/PrivateRoute";
import Register from "./Register";
import RequestPage from "./requests/Requests";
import {useAuthenticated} from "../Contexts/AuthenticationContext";
import {useWebhook} from "../Contexts/WebHookContext";
import DomainsList from "./domains/DomainsList";
import NewDomain from "./domains/NewDomain";
import ViewDomain from "./domains/ViewDomain";
import cookie from "js-cookie";


const App: FC = () => {
  const {isAuthenticated} = useAuthenticated()
  const {initialiseConnection} = useWebhook()
  React.useEffect(() => {
    if (isAuthenticated) {
      initialiseConnection()
    }
  }, [isAuthenticated])

  return (
    <div className={"app-wrapper"}>
      <Nav/>
      <Switch>
        <Route path={"/login"} exact component={() => {
          return <HandlePublicOnly><Login/></HandlePublicOnly>
        }}/>
        <Route path={"/register"} exact component={() => {
          return <HandlePublicOnly><Register/></HandlePublicOnly>
        }}/>
        <PrivateRoute path={"/domains"} exact component={DomainsList}/>
        <PrivateRoute path={"/domains/new"} exact component={NewDomain}/>
        <PrivateRoute path={"/domains/:id"} exact component={ViewDomain}/>
        <PrivateRoute path={"/stress-test"} component={RequestPage}/>
        <Redirect to={"/login"}/>
      </Switch>
    </div>
  )
}

const HandlePublicOnly: FC = (props) => {
  const history = useHistory()
  if (cookie.get("token")) {
    history.push("/stress-test")
  }
  return <>{props.children}</>
}

export default App