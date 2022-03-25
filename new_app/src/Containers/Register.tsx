import React, {FC} from "react";
import {Formik, Form, FormikHelpers} from "formik";
import {ButtonsRight, InlineLink, PageHeading, PageWrapper, Section} from "../Styled";
import FormikInput from "../Components/FormikInput";
import Button from "../Components/Button";
import * as Yup from "yup";
import useAPI from "../functions";
import FormikStandardError from "../Components/FormikStandardError";
import {useHistory} from "react-router-dom";
import {AuthWrapper} from "./Login";

interface formValues {
  email: string
  password: string
  password_confirmation: string
  name: string
}

const formSchema = Yup.object({
  "email": Yup.string().required().email().max(254).label("Email"),
  "password": Yup.string().required().min(6).max(254).label("Password"),
  "password_confirmation": Yup.string().required().min(6).max(254).label("Password confirmation").oneOf([Yup.ref('password'), null], 'Passwords must match'),
  "name": Yup.string().required().max(254).label("Name"),
})

const Register: FC = () => {
  const api = useAPI();
  const history = useHistory()
  const handleSubmit = async (values: formValues, helpers: FormikHelpers<formValues>) => {
    helpers.setStatus({})
    const res = await api.post<{ token: string, expires: number, expires_in: number }>("/auth/register", {
      email: values.email,
      name: values.name,
      password: values.password,
      password_confirmation: values.password_confirmation
    })

    if (api.handleFormikError(res, helpers)) {
      return
    }

    history.push("/login")
  }

  return (
    <PageWrapper>
      <AuthWrapper>
        <PageHeading>Register</PageHeading>
        <Formik
          validationSchema={formSchema}
          initialValues={{
            email: "",
            password: "",
            name: "",
            password_confirmation: ""
          }}
          onSubmit={handleSubmit}
        >
          {({isSubmitting}) => (
            <Section>
              <Form>
                <FormikInput name={"email"} label={"Email"}/>
                <FormikInput name={"name"} label={"Name"}/>
                <FormikInput name={"password"} label={"Password"} type={"password"}/>
                <FormikInput name={"password_confirmation"} label={"Password Confirmation"} type={"password"}/>

                <FormikStandardError/>

                <ButtonsRight spaceTop>
                  <p>Already have an account? Login <InlineLink to={"/login"}>here.</InlineLink></p>
                  <Button type={"submit"} disabled={isSubmitting} loading={isSubmitting}>
                    Register
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

export default Register