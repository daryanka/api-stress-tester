import styled from "styled-components";
import {Link} from "react-router-dom";

export const Theme = {
  black: "#272727",
  lightBlack: "#5d5d5d",
  standardShadow: "0 0 11px rgba(0,0,0,0.15)",
  standardBorderRadius: "7px",
  errorColor: "#ec5133",
  successColor: "#45d8a9",
  fontSize: "16px"
}

export type themeType = typeof Theme

export const PageWrapper = styled.div`
  margin: 0 40px;

  @media only screen and (max-width: 600px) {
    margin: 0 20px;
  }
`

export const PageHeading = styled.h1`
  font-size: 48px;
  color: ${(props) => props.theme.black};
  margin: 0 0 40px 0;
`

export const Section = styled.section`
  font-size: 28px;
  color: ${(props) => props.theme.black};
  margin-bottom: 20px;
`

export const SectionHeader = styled.h2`
  font-size: 40px;
  margin-top: 0;
`

export const ContentBox = styled.div`
  border-radius: 7px;
  box-shadow: ${props => props.theme.standardShadow};
  padding: 40px;
`

export const ButtonsRight = styled.div<{spaceTop: boolean}>`
  width: 100%;
  display: flex;
  justify-content: flex-end;
  align-items: center;
  
  button, p {
    margin-left: 20px;
  }
  
  ${props => props.spaceTop && "margin-top: 15px;"}
`

export const InlineLink = styled(Link)`
  color: deepskyblue;
  text-decoration: underline;
`