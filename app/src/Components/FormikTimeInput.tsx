import {useField, useFormikContext} from "formik";
import React, {FC, useMemo, useState} from "react";
import styled from "styled-components";
import _ from "lodash";
import {ErrorMessage, InputField, InputWrapper, Label, LabelText} from "./FormikInput";
import {SHToMinutes} from "../functions";

interface props extends React.InputHTMLAttributes<HTMLInputElement> {
  wrapperClassName?: string
  name: string
  label: string
}

const FormikInput: FC<props> = ({placeholder, wrapperClassName, className, label, ...props}) => {
  const {status} = useFormikContext()
  const [{onChange,...field}, meta] = useField(props)
  const [localVal, setLocalVal] = useState("")

  const errMsg = useMemo(() => {
    if (status && status[props.name]) {
      return <ErrorMessage>status[props.name]</ErrorMessage>
    }

    return meta.touched && meta.error && <ErrorMessage>{meta.error}</ErrorMessage>
  }, [meta, status])

  const Time = useMemo(() => {
    const {PrettyString, totalSeconds, valid} = SHToMinutes(localVal)
    if (totalSeconds > 0 && valid) {
      return PrettyString
    }
    return null
  }, [localVal])


  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    onChange(e)
    // Save locally
    setLocalVal(e.target.value)
  }

  return (
    <InputWrapper className={wrapperClassName ? wrapperClassName : ""}>
      <Label>
        <InputField
          {...props}
          onChange={handleChange}
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

export default FormikInput