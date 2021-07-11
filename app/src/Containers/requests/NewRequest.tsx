import React, {FC} from "react";
import {ButtonsRight, FormGrid, GridItem, Section, SectionHeader} from "../../Styled";
import {Form, Formik, FormikHelpers} from "formik";
import FormikSelect from "../../Components/FormikSelect";
import Button from "../../Components/Button";
import FormikTimeInput from "../../Components/FormikTimeInput";
import * as Yup from "yup";
import useAPI, {SMToMinutes} from "../../functions";
import FormikInput from "../../Components/FormikInput";
import FormikTextField from "../../Components/FormikTextField";
import useDomains from "../../Hooks/DomainsHook";
import _ from "lodash";
import FormikStandardError from "../../Components/FormikStandardError";
import {useHistory} from "react-router-dom";
import {useQueryClient} from "react-query";

interface FormValues {
  req_name: string
  method: string
  domain_id: string
  endpoint: string
  duration: string
  payload?: string
  num_requests: string
}

const validationSchema = Yup.object({
  method: Yup.string().required().label("Method"),
  req_name: Yup.string().required().label("Request Name"),
  domain_id: Yup.string().required().label("Domain"),
  endpoint: Yup.string().required().label("Endpoint"),
  num_requests: Yup.number().required().label("Number Of Requests"),
  duration: Yup.string().test({
    name: "time_type_custom",
    message: (v: { value: string | undefined }) => {
      if (v.value === undefined) {
        return "Duration format is invalid"
      }
      const {valid, totalSeconds} = SMToMinutes(v.value)
      if (valid) {
        // Check total time time is correct
        if (totalSeconds > 300) { // max 300s / 5m
          return "Maximum duration is 5 minutes"
        }
        if (totalSeconds < 10) {
          return "Minimum duration is 10 seconds"
        }
      }
      return "Duration format is invalid"
    },
    test: function (val) {
      if (val === undefined) {
        return false
      }
      const {valid, totalSeconds} = SMToMinutes(val)
      if (valid) {
        // Check total time time is correct
        if (totalSeconds > 300) { // max 300s / 5m
          return false
        }
        return totalSeconds >= 10;
      }
      return false
    }
  }),
  payload: Yup.string().label("Payload")
})

const methodOptions = [
  {
    label: "GET",
    value: "GET"
  },
  {
    label: "HEAD",
    value: "HEAD"
  },
  {
    label: "PUT",
    value: "PUT"
  },
  {
    label: "DELETE",
    value: "DELETE"
  },
  {
    label: "OPTIONS",
    value: "OPTIONS"
  },
  {
    label: "PATCH",
    value: "PATCH"
  }
]

const NewRequest: FC = () => {
  const {verifiedOnlyOption} = useDomains()
  const api = useAPI();
  const history = useHistory()
  const queryClient = useQueryClient()

  const handleSubmit = async (values: FormValues, helpers: FormikHelpers<any>) => {
    const {totalSeconds} = SMToMinutes(values.duration)
    helpers.setStatus({})

    const res = await api.post<{id: number}>("/requests/create", {
      req_name: values.req_name,
      domain_id: values.domain_id,
      time: totalSeconds,
      endpoint: values.endpoint,
      method: values.method,
      num_requests: parseInt(values.num_requests),
      payload: _.isEmpty(values.payload) ? null : values.payload
    })

    if (!api.handleFormikError(res, helpers)) {
      // Invalid List
      queryClient.invalidateQueries(["requests"])
      // Push to view page
      history.push(`/stress-test/${res.data.id}`)
    }
  }

  return (
    <>
      <SectionHeader>New Stress Test</SectionHeader>
      <Formik
        initialValues={{
          method: "",
          duration: "",
          domain_id: "",
          payload: "",
          endpoint: "",
          num_requests: "",
          req_name: "",
        }}
        validationSchema={validationSchema}
        onSubmit={handleSubmit}>
        {() => {
          return (
            <Section>
              <Form>
                <FormGrid numCols={4}>
                  <GridItem startCol={1} endCol={3}>
                    <FormikInput name={"req_name"} label={"Request Name"}/>
                  </GridItem>

                  <FormikSelect name={"method"} label={"Method"} options={methodOptions}/>
                  <FormikSelect name={"domain_id"} label={"Domain"} options={verifiedOnlyOption}/>
                  <GridItem startCol={1} endCol={5}>
                    <FormikInput name={"endpoint"} label={"Endpoint"}/>
                  </GridItem>

                  <GridItem startCol={1} endCol={3}>
                    <FormikTimeInput name={"duration"} label={"Duration (e.g. 2m 30s)"}/>
                  </GridItem>

                  <GridItem startCol={3} endCol={5}>
                    <FormikInput name={"num_requests"} label={"Number Of Requests"} type={"number"}/>
                  </GridItem>

                  <GridItem startCol={1} endCol={5}>
                    <FormikTextField name={"payload"} label={"Payload (optional)"}/>
                  </GridItem>
                </FormGrid>

                <FormikStandardError/>
                <ButtonsRight spaceTop>
                  <Button type={"submit"}>Start</Button>
                </ButtonsRight>
              </Form>
            </Section>
          )
        }}
      </Formik>
    </>
  )
}

export default NewRequest;