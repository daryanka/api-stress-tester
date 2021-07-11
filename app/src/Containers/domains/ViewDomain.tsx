import React, {FC, useMemo} from "react";
import {Prism as SyntaxHighlighter} from 'react-syntax-highlighter';
import {darcula} from 'react-syntax-highlighter/dist/esm/styles/prism';
import {RouteComponentProps, useHistory} from "react-router-dom";
import {ButtonsRight, ContentBox, PageWrapper, SectionHeader, SmallSectionHeader} from "../../Styled";
import useDomains, {domainI} from "../../Hooks/DomainsHook";
import ContentLoader from "../../Components/ContentLoader";
import styled from "styled-components";
import {Form, Formik, FormikHelpers} from "formik";
import useAPI from "../../functions";
import Button from "../../Components/Button";
import FormikInput from "../../Components/FormikInput";
import FormikStandardError from "../../Components/FormikStandardError";
import * as Yup from "yup";

const ViewDomain: FC<RouteComponentProps<{ id: string }>> = (props) => {
  const id = useMemo(() => {
    return props.match.params.id
  }, [props.match.params.id])

  const {all, query} = useDomains()
  const domain = useMemo(() => {
    return all.find((x) => `${x.id}` === id)
  }, [all])

  const refetch = () => {
    query.refetch()
  }

  return (
    <PageWrapper>
      <ContentBox>
        <SectionHeader>Domain</SectionHeader>
        <ContentLoader info={query}>
          {!domain ? (
            <p>Domain not found</p>
          ) : domain.verified === 1 ? <Verified data={domain}/> : <NotVerified refetch={refetch} data={domain}/>}
        </ContentLoader>
      </ContentBox>
    </PageWrapper>
  )
}

const VerifyDiv = styled.div`
  pre {
    border-radius: ${props => props.theme.standardBorderRadius};
  }
`

const validationSchema = Yup.object({
  endpoint: Yup.string().required().label("Endpoint")
})

const NotVerified: FC<{ data: domainI, refetch: () => void }> = (props) => {
  const api = useAPI();
  const history = useHistory()
  const jsonStr = useMemo(() => {
    return JSON.stringify({
      token: props.data.token
    }, null, "    ")
  }, [props.data])

  const fullURL = (input: string) => {
    let slash = "/"
    if (input.length > 0 && input[0] === "/") {
      slash = ""
    }
    return props.data.domain_url + slash + input
  }

  const handleSubmit = async (values: { endpoint: string }, helpers: FormikHelpers<{ endpoint: string }>) => {
    helpers.setStatus({})
    const res = await api.post("/domains/confirm", {
      endpoint: values.endpoint,
      id: props.data.id
    })

    if (!api.handleFormikError(res, helpers)) {
      props.refetch()
    }
  }

  return (
    <VerifyDiv>
      <SmallSectionHeader>Verify your domain</SmallSectionHeader>
      <p>
        Enter the URL, not including the origin, that will return the following JSON body when a GET request is made to
        it.
      </p>
      <br />
      <p>Token: {props.data.token}</p>
      <br/>
      <p>Example:</p>
      <SyntaxHighlighter language="json" style={darcula}>
        {jsonStr}
      </SyntaxHighlighter>

      <Formik
        initialValues={{endpoint: ""}}
        onSubmit={handleSubmit}
        validationSchema={validationSchema}
      >
        {({values, isSubmitting}) => {
          return (
            <Form>
              <FormikInput name={"endpoint"} label={"Endpoint"}/>

              {values.endpoint.length > 0 && <p>Full URL: {fullURL(values.endpoint)}</p>}

              <FormikStandardError spaceTop/>

              <ButtonsRight>
                <Button type={"submit"} disabled={isSubmitting} onClick={() => history.push("/domains")}>
                  Back
                </Button>
                <Button type={"submit"} disabled={isSubmitting} loading={isSubmitting}>
                  Verify
                </Button>
              </ButtonsRight>
            </Form>
          )
        }}
          </Formik>
          </VerifyDiv>
          )
        }

const Verified: FC<{ data: domainI }> = (props) => {
  const history = useHistory()
  return (
    <div>
      <p>Your domain has been verified and stress test can be made against this domain.</p>
      <ButtonsRight spaceTop>
        <Button onClick={() => history.push("/domains")}>Back</Button>
      </ButtonsRight>
    </div>
  )
}

export default ViewDomain;