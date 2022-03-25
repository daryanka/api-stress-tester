import React, {createContext, FC, useContext, useState} from "react";

interface AuthenticatedContextI {
  isAuthenticated: boolean
  setIsAuthenticated: React.Dispatch<React.SetStateAction<boolean>>
}

const Authenticated = createContext<AuthenticatedContextI | null>(null);

const AuthenticationContextProvider: FC = (props) => {
  const [isAuthenticated, setIsAuthenticated] = useState(false)

  return (
    <Authenticated.Provider value={{isAuthenticated, setIsAuthenticated}}>
      {props.children}
    </Authenticated.Provider>
  )
};

export default AuthenticationContextProvider

export const useAuthenticated = () => {
  const {isAuthenticated, setIsAuthenticated} = useContext(Authenticated) as AuthenticatedContextI;

  return {
    isAuthenticated,
    setIsAuthenticated
  }
}