import styled from "styled-components";
import React, {FC, useMemo, useState} from "react";
import useAPI from "../../functions";
import {useQuery} from "react-query";

const ReqListDiv = styled.div`
  display: flex;
  flex-direction: column;
  border-right: 2px solid ${props => props.theme.lightBlack};
`

const ReqListItem = styled.div`
  padding: 20px 5px;
  border-bottom: 2px solid ${props => props.theme.lightBlack};
  
  &:last-child {
    border-bottom: none;
  }
  
  &:hover {
    cursor: pointer;
  }
`

interface RequestData {
  id: number
  req_name: string
  user_id: number
  domain_id: number
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
}

const StandardInput = styled.input`
  border: 1px solid ${props => props.theme.lightBlack};
  outline: none;
  padding: 10px;
  border-radius: 7px;
  width: 95%;
`

const RequestsList: FC = () => {
  const [search, setSearch] = useState("")
  const api = useAPI()
  const data = useQuery(["requests"], async () => {
    const res = await api.get<RequestData[]>("/requests/all")
    if (!api.error(res)) {
      return res.data
    }
  })

  const DataList = useMemo(() => {
    if (data.data) {
      // Check if it should filter by search
      if (search !== "") {
        const newData = [];

        for (let i = 0; i < data.data.length; i++) {
          const curr = data.data[i].req_name.toLowerCase()
          if (curr.includes(search.toLowerCase())) {
            // Add
            newData.push({...data.data[i]})
          }
        }

        return newData
      }

      return data.data
    }

    return []
  }, [data.data, search])

  return (
    <ReqListDiv>
      <label>
        <p style={{marginBottom: "10px"}}>Search</p>
        <StandardInput value={search} onChange={e => setSearch(e.target.value)} />
      </label>
      {DataList.map((el) => {
        return (
          <ReqListItem key={`request-listitem-${el.id}`}>
            Request: {el.req_name}
          </ReqListItem>
        )
      })}
    </ReqListDiv>
  )
}

export default RequestsList