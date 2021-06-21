import React, {FC, useMemo, useState} from "react";
import {ContentBox, PageWrapper} from "../../Styled";
import styled from "styled-components";
import RequestsList from "./ListRequests";
import {Route, RouteChildrenProps, Switch} from "react-router-dom";
import NewRequest from "./NewRequest";
import RequestView from "./RequestView";

const Split = styled.div`
  display: grid;
  grid-template-columns: 200px auto;
`

const RequestPage: FC<RouteChildrenProps> = (props) => {
  return (
    <PageWrapper>
      <ContentBox>
        <Split>
          <RequestsList/>
          <RequestData {...props} />
        </Split>
      </ContentBox>
    </PageWrapper>
  )
}

const RequestViewWrapper = styled.div`
  padding: 0 40px 40px 40px;
`

const RequestData: FC<RouteChildrenProps> = (props) => {
  return (
    <RequestViewWrapper>
      <Switch>
        <Route path={props.match?.path} exact render={() => {
          return (
            <Temp name={"homepage"}/>
          )
        }
        }/>
        <Route path={`${props.match?.path}/create`} component={NewRequest} />
        <Route path={`${props.match?.path}/:id`} exact component={RequestView}/>
      </Switch>
    </RequestViewWrapper>
  )
}

const Temp: FC<{ name: string }> = (props) => {
  return (
    <h1>{props.name}</h1>
  )
}

export default RequestPage