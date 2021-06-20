import React, {FC} from "react";
import {ButtonsRight, FormGrid, GridItem, Section, SectionHeader} from "../../Styled";
import {Form, Formik} from "formik";
import FormikSelect from "../../Components/FormikSelect";
import Button from "../../Components/Button";
import FormikTimeInput from "../../Components/FormikTimeInput";
import * as Yup from "yup";
import {SMToMinutes} from "../../functions";
import FormikInput from "../../Components/FormikInput";
import FormikTextField from "../../Components/FormikTextField";
import useDomains from "../../Hooks/DomainsHook";

interface FormValues {
  method: string
  domain: string
  endpoint: string
  duration: string
  payload?: string
}

const validationSchema = Yup.object({
  method: Yup.string().required().label("Method"),
  domain: Yup.string().required().label("Domain"),
  endpoint: Yup.string().required().label("Endpoint"),
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

  const test = (values: FormValues) => {
    console.log(values)
  }

  return (
    <>
      <SectionHeader>New Stress Test</SectionHeader>
      <Formik
        initialValues={{
          method: "",
          duration: "",
          domain: "",
          payload: "",
          endpoint: ""
        }}
        validationSchema={validationSchema}
        onSubmit={test}>
        {() => {
          return (
            <Section>
              <Form>
                <FormGrid numCols={3}>
                  <FormikSelect name={"method"} label={"Method"} options={methodOptions}/>
                  <FormikTimeInput name={"duration"} label={"Duration (e.g. 2m 30s)"}/>
                  <FormikSelect name={"domain"} label={"Domain"} options={verifiedOnlyOption}/>
                  <GridItem startCol={1} endCol={4}>
                    <FormikInput name={"endpoint"} label={"Endpoint"}/>
                  </GridItem>
                  <GridItem startCol={1} endCol={4}>
                    <FormikTextField name={"payload"} label={"Payload (optional)"}/>
                  </GridItem>
                </FormGrid>
                <ButtonsRight spaceTop>
                  <Button>Start</Button>
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