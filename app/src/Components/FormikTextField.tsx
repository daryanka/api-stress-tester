import {useField, useFormikContext} from "formik";
import React, {FC, useMemo} from "react";
import _ from "lodash";
import {ErrorMessage, InputField, InputWrapper, Label, LabelText} from "./FormikInput";
import styled from "styled-components";

interface props extends React.InputHTMLAttributes<HTMLTextAreaElement> {
  wrapperClassName?: string
  name: string
  label: string
}

const TextField = styled(InputField)`
  min-height: 100px;
  max-height: 500px;
  resize: vertical;
`

const FormikTextField: FC<props> = ({placeholder, wrapperClassName, className, label, ...props}) => {
  const {status} = useFormikContext()
  const [field, meta] = useField(props)

  const errMsg = useMemo(() => {
    if (status && status[props.name]) {
      return <ErrorMessage>status[props.name]</ErrorMessage>
    }

    return meta.touched && meta.error && <ErrorMessage>{meta.error}</ErrorMessage>
  }, [meta, status])

  return (
    <InputWrapper className={wrapperClassName ? wrapperClassName : ""}>
      <Label>
        <TextField
          as={"textarea"}
          {...props}
          {...field}
          className={`${className} ${errMsg && "error"} ${meta.touched && "touched"} ${!_.isEmpty(field.value) && "has-val"}`}
        />
        {label && <LabelText className={"label-holder"}>{label}</LabelText>}
        {errMsg}
      </Label>
    </InputWrapper>
  )
}

export default FormikTextField