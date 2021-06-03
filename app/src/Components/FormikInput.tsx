import {useField, useFormikContext} from "formik";
import React, {FC, useMemo} from "react";
import styled from "styled-components";

interface props extends React.InputHTMLAttributes<HTMLInputElement> {
  wrapperClassName?: string
  name: string
  label: string
}

const InputField = styled.input`

`

const ErrorMessage = styled.p`
`

const Label = styled.p``

const FormikInput: FC<props> = ({wrapperClassName, className, label, ...props}) => {
  const {status} = useFormikContext()
  const [field, meta] = useField(props)

  const errMsg = useMemo(() => {
    if (status[props.name]) {
      return <ErrorMessage>status[props.name]</ErrorMessage>
    }

    return meta.touched && meta.error && <ErrorMessage>{meta.error}</ErrorMessage>
  }, [meta, status])

  return (
    <div className={wrapperClassName ? wrapperClassName : ""}>
      <label>
        {label && <Label>{label}</Label>}
        <InputField
          {...props}
          {...field}
          className={`${className} ${errMsg ? "error" : ""} ${meta.touched ? "touched" : ""}`}
        />
        {errMsg}
      </label>
    </div>
  )
}