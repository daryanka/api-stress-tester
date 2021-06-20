import {Route, useHistory} from "react-router-dom";
import {useQuery} from "react-query";
import {FC} from "react";
import useAPI from "../functions";
import React from "react";
import {RouteProps} from "react-router-dom";
import {useAuthenticated} from "../Contexts/AuthenticationContext";

const PrivateRoute: FC<RouteProps> = (props) => {
  const {setIsAuthenticated} = useAuthenticated()
  const api = useAPI()
  const history = useHistory()
  const me = useQuery(["me"], async () => {
    const res = await api.get("/auth/me")
    if (res.status !== 200) {
      throw Error("Unauthorised")
    }
  }, {
    onError: () => {
      api.logout()
    },
    onSuccess: () => {
      setIsAuthenticated(true)
    },
    refetchOnWindowFocus: true,
  })

  if (me.isLoading) {
    return <h1>Loading</h1>
  }
  return <Route {...props} />
}

export default PrivateRoute