import React, {FC} from "react";
import ReactDOM from "react-dom";
import "./styles.scss";
import App from "./Containers/App";
import {BrowserRouter as Router} from "react-router-dom";
import {WebhookProvider} from "./Contexts/WebHookContext";
import styled, {ThemeProvider} from "styled-components";
import {Theme} from "./Styled";
import {QueryClientProvider, QueryClient} from "react-query";
import AuthenticationContext from "./Contexts/AuthenticationContext";

const General = styled.div`
  * {
    font-family: Roboto, serif;
  }

  p {
    font-size: ${props => props.theme.fontSize};
  }
`

const client = new QueryClient({
  defaultOptions: {
    queries: {
      retry: false,
      refetchOnWindowFocus: false
    }
  }
})

const RootApp: FC = () => {
  return (
    <ThemeProvider theme={Theme}>
      <QueryClientProvider client={client}>
        <AuthenticationContext>
          <General>
            <WebhookProvider>
              <Router>
                <App/>
              </Router>
            </WebhookProvider>
          </General>
        </AuthenticationContext>
      </QueryClientProvider>
    </ThemeProvider>
  )
}

ReactDOM.render(<RootApp/>, document.getElementById("root"))