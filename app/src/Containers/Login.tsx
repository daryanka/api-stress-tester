import React, {FC} from "react";
import {Formik, Form, FormikHelpers} from "formik";
import {ButtonsRight, InlineLink, PageHeading, PageWrapper, Section} from "../Styled";
import FormikInput from "../Components/FormikInput";
import Button from "../Components/Button";
import * as Yup from "yup";
import styled from "styled-components";
import useAPI from "../functions";
import FormikStandardError from "../Components/FormikStandardError";
import cookie from "js-cookie";
import {useHistory} from "react-router-dom";

interface formValues {
  email: string
  password: string
}

const formSchema = Yup.object({
  "email": Yup.string().required().email().label("Email"),
  "password": Yup.string().required().label("Password"),
})

export const AuthWrapper = styled.div`
  max-width: 800px;
  width: 80%;
  margin: auto;
  padding: 40px;
  box-shadow: ${props => props.theme.standardShadow};
  border-radius: ${props => props.theme.standardBorderRadius};
`

const Login: FC = () => {
  const api = useAPI();
  const history = useHistory()
  const handleSubmit = async (values: formValues, helpers: FormikHelpers<formValues>) => {
    helpers.setStatus({})
    const res = await api.post<{token: string, expires: number, expires_in: number}>("/auth/login", {
      email: values.email,
      password: values.password
    })

    if (api.handleFormikError(res, helpers)) {
      helpers.setTouched({})
      return
    }

    // Save token and redirect to dashboard
    cookie.set("token", res.data.token, {
      expires: new Date(res.data.expires * 1000)
    })

    history.replace("/dashboard")
  }

  return (
    <PageWrapper>
      <AuthWrapper>
        <PageHeading>Login</PageHeading>
        <Formik
          validationSchema={formSchema}
          initialValues={{
            email: "",
            password: ""
          }}
          onSubmit={handleSubmit}
        >
          {({isSubmitting}) => (
            <Section>
              <Form>
                <FormikInput name={"email"} label={"Email"}/>
                <FormikInput name={"password"} label={"Password"} type={"password"}/>

                <FormikStandardError />

                <ButtonsRight spaceTop>
                  <p>Don't have an account? Register <InlineLink to={"/register"}>here.</InlineLink></p>
                  <Button type={"submit"} disabled={isSubmitting} loading={isSubmitting}>
                    Login
                  </Button>
                </ButtonsRight>
              </Form>
            </Section>
          )}
        </Formik>
      </AuthWrapper>
    </PageWrapper>
  )
}

export default Login