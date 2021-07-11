import React, {FC} from "react";
import {ButtonsRight, ContentBox, PageWrapper, Section, SectionHeader, SmallSectionHeader} from "../../Styled";
import * as Yup from "yup";
import {Formik, Form, FormikHelpers} from "formik";
import FormikInput from "../../Components/FormikInput";
import useAPI from "../../functions";
import Button from "../../Components/Button";
import FormikStandardError from "../../Components/FormikStandardError";

const validURL = (str: string) => {
  const pattern = new RegExp('^(https?:\\/\\/)?' + // protocol
    '((([a-z\\d]([a-z\\d-]*[a-z\\d])*)\\.)+[a-z]{2,}|' + // domain name
    '((\\d{1,3}\\.){3}\\d{1,3}))' + // OR ip (v4) address
    '(\\:\\d+)?(\\/[-a-z\\d%_.~+]*)*' + // port and path
    '(\\?[;&a-z\\d%_.~+=-]*)?' + // query string
    '(\\#[-a-z\\d_]*)?$', 'i'); // fragment locator
  return pattern.test(str);
}


const formSchema = Yup.object({
  domain_url: Yup.string().max(254).required().test({
    name: "url_type",
    message: (v: { value: string | undefined }) => {
      if (v.value === undefined) {
        return "Domain URL is a required field"
      }
      return "Domain URL format is invalid"
    },
    test: function (val) {
      if (val === undefined) {
        return false
      }
      try {
        const r = new URL(val)
        if (r.origin === "null") {
          return false
        }
        return validURL(val)
      } catch (e) {
        return false
      }
    }
  }).label("Domain URL"),
})

const NewDomain: FC = () => {
  const api = useAPI()
  const handleSubmit = async (values: { domain_url: string }, helpers: FormikHelpers<{ domain_url: string }>) => {
    helpers.setStatus({})
    let origin = ""
    try {
      const urlObj = new URL(values.domain_url)
      const domainURL = urlObj.origin
      if (domainURL === "null") {
        throw Error("Invalid domain")
      }
      origin = domainURL
    } catch (e) {
      helpers.setStatus({
        domain_url: "Domain URL is a required field"
      })
      return
    }


    const res = await api.post<{id: number}>("/v1/domains/create", {
      domain_url: origin
    })

    if (!api.handleFormikError(res, helpers)) {

    }
  }

  const domainValue = (val: string): string => {
    try {
      const urlVal = new URL(val)
      if (!validURL(val)) {
        return ""
      }
      return urlVal.origin === "null" ? "" : urlVal.origin
    } catch (e) {
      return ""
    }
  }

  return (
    <PageWrapper>
      <ContentBox>
        <SectionHeader>New Domain</SectionHeader>
        <SmallSectionHeader>How it works</SmallSectionHeader>
        <p>
          Before testing any API we first need to verify that the domain being testing belongs to you.
          To verify the domain belongs to you there are 3 simple steps:

          <br/>
          <br/>
          <strong>1.</strong> Register your full domain URL.
          <br/>
          <br/>
          <strong>2.</strong> Add a new endpoint your on API which returns a token you will be provided with.
          <br/>
          <br/>
          <strong>3.</strong> Confirm the endpoint is registered and let us check against your endpoint that the valid
          token is being returned.
          <br/>
          <br/>
        </p>

        <Section
          spaceTop
        >
          <Formik
            onSubmit={handleSubmit}
            initialValues={{domain_url: ""}}
            validationSchema={formSchema}
          >
            {({values, isSubmitting}) => (
              <Form>
                <p>The domain URL should include the protocol and port if the port is not 80.</p>
                <FormikInput name={"domain_url"} label={"Domain URL"}/>
                <p>{domainValue(values.domain_url)}</p>

                <FormikStandardError spaceTop/>

                <ButtonsRight>
                  <Button type={"submit"} disabled={isSubmitting} loading={isSubmitting}>Add Domain</Button>
                </ButtonsRight>
              </Form>
            )}
          </Formik>
        </Section>
      </ContentBox>
    </PageWrapper>
  )
}

// const DomainProgressBar: FC = () => {
//   return (
//     <div>
//
//     </div>
//   )
// }

export default NewDomain;