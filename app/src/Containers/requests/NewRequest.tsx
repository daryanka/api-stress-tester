import React, {FC} from "react";
import {ButtonsRight, SectionHeader} from "../../Styled";
import {Form, Formik} from "formik";
import FormikSelect from "../../Components/FormikSelect";
import Button from "../../Components/Button";

interface FormValues {
  method: string
}

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
  const test = (values: FormValues) => {
    console.log(values)
  }

  return (
    <>
      <SectionHeader>New Stress Test</SectionHeader>
      <Formik
        initialValues={{
          method: ""
        }}
        onSubmit={test}>
        {() => {
          return (
            <Form>
              <FormikSelect name={"method"} label={"Method"} options={methodOptions}/>
              <ButtonsRight spaceTop>
                <Button>Start</Button>
              </ButtonsRight>
            </Form>
          )
        }}
      </Formik>
    </>
  )
}

export default NewRequest;