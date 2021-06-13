import React, {FC, useMemo} from "react";
import styled from "styled-components";
import {PageWrapper} from "../Styled";
import {Link} from "react-router-dom";
import {useQuery} from "react-query";
import {useAuthenticated} from "../Contexts/AuthenticationContext";
import useAPI from "../functions";

const NavBar = styled.nav`
  height: 70px;
  width: 100%;
  margin-bottom: 70px;
  display: flex;
  justify-content: space-between;
  align-items: center;
`

const List = styled.ul`
  display: flex;
  list-style: none;
`

const ListItem = styled(Link)`
  margin: 0 20px;
  color: ${props => props.theme.black};
  text-decoration: none;
  font-size: 20px;
`

const LogoutItem = styled.span`
  margin: 0 20px;
  color: ${props => props.theme.black};
  text-decoration: none;
  font-size: 20px;
  
  &:hover {
    cursor: pointer;
  }
`

const Nav: FC = () => {
  const {logout} = useAPI()
  const {isAuthenticated} = useAuthenticated();

  return (
    <PageWrapper>
      <NavBar>
        <h3>API Stress Tester</h3>
        <List>
          {isAuthenticated ? (
            <>
              <ListItem to={"/dashboard"}>Dashboard</ListItem>
              <ListItem to={"/domains"}>Domains</ListItem>
              <ListItem to={"/stress-test"}>Stress Tests</ListItem>
              <LogoutItem onClick={logout}>Logout</LogoutItem>
            </>
          ) : (
            <>
              <ListItem to={"/"}>Home</ListItem>
              <ListItem to={"/solution"}>Solution</ListItem>
              <ListItem to={"/login"}>Login</ListItem>
              <ListItem to={"/register"}>Register</ListItem>
            </>
          )}
        </List>
      </NavBar>
    </PageWrapper>
  )
}

export default Nav
