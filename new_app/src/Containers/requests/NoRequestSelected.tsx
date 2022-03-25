import React, {FC, useMemo} from "react";
import {SectionHeader, SmallSectionHeader} from "../../Styled";
import ContentLoader from "../../Components/ContentLoader";
import useRequestList from "../../Hooks/RequestListHook";
import styled from "styled-components";
import Button from "../../Components/Button";
import _ from "lodash";
import {RequestData} from "./ListRequests";
import {useHistory} from "react-router-dom";

const Boxes = styled.div`
  display: grid;
  grid-template-columns: 1fr 1fr;
  grid-gap: 20px;
`

const Box = styled.div<{ withBorder?: boolean }>`
  ${props => props.withBorder && `border-right: 2px solid ${props.theme.lightBlack};`}
`

const Center = styled.div`
  display: flex;
  width: 100%;
  margin-top: 20px;
  justify-content: center;
`

const Recent = styled.div<{last?: boolean}>`
  padding: 10px 10px 20px 10px;
  margin: 10px;
  border-bottom: 2px solid ${props => props.last ? "white": props.theme.lightBlack};
  
  p {
    margin-bottom: 10px;
  }
  
  .green {
    color: ${props => props.theme.successColor};
  }

  .red {
    color: ${props => props.theme.errorColor};
  }
  
  .b {
    font-weight: bold;
  }

`


const RecentRequest: FC<{data: RequestData, last: boolean}> = ({data, last}) => {
  const history = useHistory()

  return (
    <Recent last={last}>
      <p className={"b"}>{data.req_name}:</p>
      <p className={"green"}>{data.successful_req} Successful Requests</p>
      <p className={"red"}>{data.failed_req} Failed Requests</p>
      <p>{data.average_response_time.toFixed(2)}ms Average Response Time</p>
      <Center><Button onClick={() => history.push(`/stress-test/${data.id}`)}>View</Button></Center>
    </Recent>
  )
}

const NoRequestSelected: FC = () => {
  const history = useHistory()
  const requestList = useRequestList()
  const recentRequest = useMemo(() => {
    // Get first 5 requests
    const res = [];
    if (!_.isEmpty(requestList.data)) {
      for (let i = 0; i < requestList.data!.length; i++) {
        res.push({...requestList.data![i]})
        if (i == 4) {
          break
        }
      }
    }
    return res
  }, [requestList.data])

  return (
    <ContentLoader info={requestList}>
      <SectionHeader>Stress Tests</SectionHeader>
      <Boxes>
        <Box withBorder>
          <SmallSectionHeader center>Recent Stress Tests</SmallSectionHeader>
          {recentRequest.length === 0 ? (
            <div>
              No stress tests found
            </div>
          ) : recentRequest.map((el, index) => {
            return (
              <RecentRequest key={`recent-req-${el.id}`} data={el} last={index === recentRequest.length - 1} />
            )
          })}
        </Box>
        <Box>
          <SmallSectionHeader center>New Stress Test</SmallSectionHeader>
          <Center>
            <Button onClick={() => history.push(`/stress-test/create`)}>New Stress Test</Button>
          </Center>
        </Box>
      </Boxes>
    </ContentLoader>
  )
}

export default NoRequestSelected;