import React, {FC} from "react";
import styled from "styled-components";
import {Query, UseQueryResult} from "react-query";

const Spinner = styled.div`
  position: absolute;
  width: 100px;
  height: 100px;

  .sk-circle {
    width: 100%;
    height: 100%;
    position: absolute;
    left: 0;
    top: 0;
  }

  .sk-circle:before {
    content: '';
    display: block;
    margin: 0 auto;
    width: 15%;
    height: 15%;
    background-color: ${props => props.theme.lightBlack};
    border-radius: 100%;
    -webkit-animation: sk-circleFadeDelay 1.2s infinite ease-in-out both;
    animation: sk-circleFadeDelay 1.2s infinite ease-in-out both;
  }

  .sk-circle2 {
    -webkit-transform: rotate(30deg);
    -ms-transform: rotate(30deg);
    transform: rotate(30deg);
  }

  .sk-circle3 {
    -webkit-transform: rotate(60deg);
    -ms-transform: rotate(60deg);
    transform: rotate(60deg);
  }

  .sk-circle4 {
    -webkit-transform: rotate(90deg);
    -ms-transform: rotate(90deg);
    transform: rotate(90deg);
  }

  .sk-circle5 {
    -webkit-transform: rotate(120deg);
    -ms-transform: rotate(120deg);
    transform: rotate(120deg);
  }

  .sk-circle6 {
    -webkit-transform: rotate(150deg);
    -ms-transform: rotate(150deg);
    transform: rotate(150deg);
  }

  .sk-circle7 {
    -webkit-transform: rotate(180deg);
    -ms-transform: rotate(180deg);
    transform: rotate(180deg);
  }

  .sk-circle8 {
    -webkit-transform: rotate(210deg);
    -ms-transform: rotate(210deg);
    transform: rotate(210deg);
  }

  .sk-circle9 {
    -webkit-transform: rotate(240deg);
    -ms-transform: rotate(240deg);
    transform: rotate(240deg);
  }

  .sk-circle10 {
    -webkit-transform: rotate(270deg);
    -ms-transform: rotate(270deg);
    transform: rotate(270deg);
  }

  .sk-circle11 {
    -webkit-transform: rotate(300deg);
    -ms-transform: rotate(300deg);
    transform: rotate(300deg);
  }

  .sk-circle12 {
    -webkit-transform: rotate(330deg);
    -ms-transform: rotate(330deg);
    transform: rotate(330deg);
  }

  .sk-circle2:before {
    -webkit-animation-delay: -1.1s;
    animation-delay: -1.1s;
  }

  .sk-circle3:before {
    -webkit-animation-delay: -1s;
    animation-delay: -1s;
  }

  .sk-circle4:before {
    -webkit-animation-delay: -0.9s;
    animation-delay: -0.9s;
  }

  .sk-circle5:before {
    -webkit-animation-delay: -0.8s;
    animation-delay: -0.8s;
  }

  .sk-circle6:before {
    -webkit-animation-delay: -0.7s;
    animation-delay: -0.7s;
  }

  .sk-circle7:before {
    -webkit-animation-delay: -0.6s;
    animation-delay: -0.6s;
  }

  .sk-circle8:before {
    -webkit-animation-delay: -0.5s;
    animation-delay: -0.5s;
  }

  .sk-circle9:before {
    -webkit-animation-delay: -0.4s;
    animation-delay: -0.4s;
  }

  .sk-circle10:before {
    -webkit-animation-delay: -0.3s;
    animation-delay: -0.3s;
  }

  .sk-circle11:before {
    -webkit-animation-delay: -0.2s;
    animation-delay: -0.2s;
  }

  .sk-circle12:before {
    -webkit-animation-delay: -0.1s;
    animation-delay: -0.1s;
  }

  @-webkit-keyframes sk-circleFadeDelay {
    0%, 39%, 100% {
      opacity: 0;
    }
    40% {
      opacity: 1;
    }
  }

  @keyframes sk-circleFadeDelay {
    0%, 39%, 100% {
      opacity: 0;
    }
    40% {
      opacity: 1;
    }
  }
`

const ContentLoaderDiv = styled.div`
  margin-top: 50px;
  height: 100%;
  align-items: center;
  display: flex;
  justify-content: center;
`

const ContentLoader: FC<{info: UseQueryResult}> = (props) => {
  if (props.info.isLoading) {
    return (
      <div>
        <ContentLoaderDiv>
          <Spinner>
            <div className="sk-circle1 sk-circle"/>
            <div className="sk-circle2 sk-circle"/>
            <div className="sk-circle3 sk-circle"/>
            <div className="sk-circle4 sk-circle"/>
            <div className="sk-circle5 sk-circle"/>
            <div className="sk-circle6 sk-circle"/>
            <div className="sk-circle7 sk-circle"/>
            <div className="sk-circle8 sk-circle"/>
            <div className="sk-circle9 sk-circle"/>
            <div className="sk-circle10 sk-circle"/>
            <div className="sk-circle11 sk-circle"/>
            <div className="sk-circle12 sk-circle"/>
          </Spinner>
        </ContentLoaderDiv>
      </div>
    )
  }

  if (props.info.error) {
    let message = "Something went wrong please try again later."
    if (typeof props.info.error === "string") {
      message = props.info.error
    }
    return (
      <p>
        {message}
      </p>
    )
  }

  return <>{props.children}</>
}

export default ContentLoader;