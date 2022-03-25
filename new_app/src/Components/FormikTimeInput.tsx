import {useField, useFormikContext} from "formik";
import React, {FC, useMemo, useState} from "react";
import styled from "styled-components";
import _ from "lodash";
import {ErrorMessage, InputField, InputWrapper, Label, LabelText} from "./FormikInput";
import {SMToMinutes} from "../functions";

interface props extends React.InputHTMLAttributes<HTMLInputElement> {
  wrapperClassName?: string
  name: string
  label: string
}

const PrettyStr = styled.p`
  font-size: 18px;
  margin-left: 10px;
  margin-top: -10px;
`

const FormikTimeInput: FC<props> = ({placeholder, wrapperClassName, className, label, ...props}) => {
  const {status} = useFormikContext()
  const [{onChange,...field}, meta] = useField(props)

  const errMsg = useMemo(() => {
    if (status && status[props.name]) {
      return <ErrorMessage>status[props.name]</ErrorMessage>
    }

    return meta.touched && meta.error && <ErrorMessage>{meta.error}</ErrorMessage>
  }, [meta, status])

  const Time = useMemo(() => {
    const {PrettyString, totalSeconds, valid} = SMToMinutes(field.value)
    if (totalSeconds > 0 && valid) {
      return <PrettyStr>{PrettyString}</PrettyStr>
    }
    return null
  }, [field.value])

  return (
    <InputWrapper className={wrapperClassName ? wrapperClassName : ""}>
      <Label>
        <InputField
          {...props}
          onChange={onChange}
          {...field}
          className={`${className} ${errMsg && "error"} ${meta.touched && "touched"} ${!_.isEmpty(field.value) && "has-val"}`}
        />
        {label && <LabelText className={"label-holder"}>{label}</LabelText>}
        {Time}
        {errMsg}
      </Label>
    </InputWrapper>
  )
}

export default FormikTimeInput