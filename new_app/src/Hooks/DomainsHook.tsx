import {useQuery} from "react-query";
import useAPI from "../functions";
import {useMemo} from "react";

export interface domainI {
  id: number
  domain_url: string
  user_id: string,
  verified: number,
  token: string,
  num_requests: number
}

const useDomains = () => {
  const api = useAPI()
  const queryInfo = useQuery(["user-domains"], async () => {
    const res = await api.get<domainI[]>("/domains/all")
    const err = api.error(res)
    if (!err) {
      return res.data
    } else {
      throw err.error
    }
  }, {
    refetchOnWindowFocus: true,
  })

  const verifiedOnly = useMemo(() => {
    if (queryInfo.data) {
      return queryInfo.data.filter(x => x.verified === 1)
    }
    return []
  }, [queryInfo.data])

  const all = useMemo(() => {
    if (queryInfo.data) {
      return queryInfo.data
    }
    return []
  }, [queryInfo.data])

  const verifiedOnlyOption = useMemo(() => {
    if (queryInfo.data) {
      return queryInfo.data.filter(x => x.verified === 1).map(el => ({
        value: el.id,
        label: el.domain_url
      }))
    }
    return []
  }, [queryInfo.data])

  return {
    verifiedOnly,
    all,
    query: queryInfo,
    verifiedOnlyOption
  }
}

export default useDomains