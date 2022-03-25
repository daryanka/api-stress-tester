import {useQuery} from "react-query";
import {RequestData} from "../Containers/requests/ListRequests";
import useAPI from "../functions";

const useRequestList = () => {
  const api = useAPI();
  return useQuery(["requests"], async () => {
    const res = await api.get<RequestData[]>("/requests/all")
    if (!api.error(res)) {
      return res.data
    } else if (res.status === 404) {
      return []
    }
  })
}

export default useRequestList