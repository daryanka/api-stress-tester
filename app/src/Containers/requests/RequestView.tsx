import React, {FC, useMemo} from "react";
import {RouteComponentProps} from "react-router-dom";
import {Section, SectionHeader, SmallSectionHeader, themeType} from "../../Styled";
import {useQuery} from "react-query";
import useAPI, {SecondsToPretty} from "../../functions";
import ContentLoader from "../../Components/ContentLoader";
import styled, {useTheme} from "styled-components";
import {ResponsiveLineCanvas, PointTooltip} from "@nivo/line";
import _ from "lodash";
import {useRequestHook} from "../../Contexts/WebHookContext";
import Button from "../../Components/Button";

interface propsI extends RouteComponentProps<{ id: string }> {}

interface RequestData {
  id: number
  req_name: string
  user_id: number
  domain_id: number
  domain: {
    id: number
    domain_url: string
  }
  endpoint: string
  method: string
  payload?: string
  time: number
  num_requests: number
  successful_req: number
  failed_req: number
  average_response_time: number
  status: number
  created_at: string
  individual_requests: {
    id: number
    created_at: string
    request_overview_id: number
    status_code: number
    time_taken: number
  }[]
}

const EndpointWrapper = styled.div`
  overflow-x: auto;
  display: grid;
  grid-template-columns: 1fr;
  padding: 10px;
  background: ${props => props.theme.grey};
  border-radius: ${props => props.theme.standardBorderRadius};
  font-weight: bold;
  margin-bottom: 30px;

  .full-endpoint {
    word-break: break-all;
  }

  .endpoint {
    display: flex;
    flex-direction: row;
  }

  .payload {
    margin-top: 20px;
  }
`

const Divider = styled.div`
  min-height: 100%;
  width: 2px;
  margin: 0 15px;
  background: ${props => props.theme.black};
`

const Pre = styled.pre`
  font-size: 16px;
`

const RequestView: FC<propsI> = (props) => {
  const {fullRequest, sendCancel} = useRequestHook(props.match.params.id)
  const theme = useTheme() as themeType
  const api = useAPI()

  const reqInfo = useQuery(["individual-request", props.match.params.id], async () => {
    const res = await api.get<RequestData>(`/requests/individual/${props.match.params.id}`)
    const err = api.error(res)
    if (!err) {
      return res.data
    } else {
      throw Error("")
    }
  })

  const chartData = useMemo(() => {
    const d = []
    if (!_.isEmpty(reqInfo.data?.individual_requests)) {
      for (let i = 0; i < reqInfo.data!.individual_requests.length; i++) {
        d.push({
          y: reqInfo.data!.individual_requests[i].time_taken,
          x: i
        })
      }
    }

    if (fullRequest?.reqs) {
      for (let i = 0; i < fullRequest?.reqs.length; i++) {
        d.push({
          x: d.length,
          y: fullRequest.reqs[i].time_taken
        })
      }
    }
    return d
  }, [fullRequest?.reqs, reqInfo.data])

  const payloadJSON = useMemo(() => {
    if (reqInfo.data?.payload) {
      return (
        <Pre>
         {JSON.stringify(JSON.parse(reqInfo.data.payload), null, "    ")}
       </Pre>
      )
    }
    return <Pre>No Payload provided</Pre>
  }, [reqInfo.data])

  const fullEndpoint = useMemo(() => {
    if (reqInfo.data) {
      let between = ""
      if (reqInfo.data.endpoint.length > 0 && reqInfo.data.endpoint[0] !== "/") {
        between = "/"
      }
      return reqInfo.data.domain.domain_url + between + reqInfo.data.endpoint
    }

    return ""
  }, [reqInfo.data])

  const prettyTime = useMemo(() => {
    if (reqInfo.data?.time) {
      return SecondsToPretty(reqInfo.data?.time)
    }
    return ""
  }, [reqInfo.data])

  const requestStatus = useMemo(() => {
    switch (reqInfo.data?.status) {
      case 1:
        return "In Progress"
      case 2:
        return "Completed"
      case 3:
        return "Cancelled"
      default:
        return ""
    }
  }, [reqInfo.data])

  const handleCancel = () => {
    sendCancel()
    setTimeout(() => {
      reqInfo.refetch()
    }, 100)
  }

  return (
    <ContentLoader info={reqInfo}>
      <SectionHeader>
        {reqInfo.data?.req_name}
      </SectionHeader>

      <Section>
        <SmallSectionHeader>Request Information</SmallSectionHeader>
        <EndpointWrapper>
          <div className={"endpoint"}>
            <p>{reqInfo.data?.method}</p>
            <Divider/>
            <p className={"full-endpoint"}>{fullEndpoint}</p>
            <Divider/>
            <p>{reqInfo.data?.num_requests} Requests</p>
            <Divider/>
            <p>{prettyTime}</p>
          </div>
          <div className="payload">
            <p>Payload:</p>
            {payloadJSON}
          </div>
        </EndpointWrapper>

        <SmallSectionHeader>Status: {requestStatus}</SmallSectionHeader>
        {/* Allow Cancel if Request is progress */}
        {reqInfo.data?.status === 1 && (
          <Button disabled={reqInfo.isLoading || reqInfo.isFetching} onClick={handleCancel}>Cancel</Button>
        )}

        <div
          style={{
            height: "600px"
          }}
        >
          <ResponsiveLineCanvas
            data={[
              {
                id: "data",
                data: chartData
              }
            ]}
            margin={{top: 20, right: 20, bottom: 50, left: 50}}
            xScale={{type: "linear"}}
            yFormat=" >-.2f"
            curve="linear"
            axisBottom={{
              legend: "Request",
              legendOffset: 36,
              legendPosition: "middle"
            }}
            axisLeft={{
              legend: "Time (ms)",
              legendOffset: -40,
              legendPosition: "middle"
            }}
            enableGridX={false}
            colors={theme.lightBlack}
            lineWidth={1}
            pointSize={0}
            pointColor={{theme: "background"}}
            pointBorderWidth={1}
            tooltip={ToolTip}
          />
        </div>
      </Section>

    </ContentLoader>
  )
}

const ToolTip: PointTooltip = ({point}) => {
  return <p>{point.data.y}ms</p>
}

export default RequestView;