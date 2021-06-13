import {useField, useFormikContext} from "formik";
import React, {FC, useMemo} from "react";
import styled from "styled-components";
import _ from "lodash";

interface props extends React.InputHTMLAttributes<HTMLInputElement> {
  wrapperClassName?: string
  name: string
  label: string
}

const InputField = styled.input`
  width: 100%;
  border: 3px solid ${props => props.theme.black};
  border-top: 0;
  border-right: 0;
  border-left: 0;
  padding: 10px;
  transition: ease all 200ms;
  margin-bottom: 20px;

  outline: none;

  &.touched {
    border-color: ${props => props.theme.successColor};

    + .label-holder {
      color: ${props => props.theme.successColor};
    }

    &.error {
      border-color: ${props => props.theme.errorColor};
      color: ${props => props.theme.errorColor};

      + .label-holder {
        color: ${props => props.theme.errorColor};
      }
    }
  }

  font-size: 20px;

  &:focus + .label-holder {
    top: -20px;
    font-size: ${props => props.theme.fontSize};
    font-weight: bold;
  }

  &.has-val + .label-holder {
    top: -20px;
    font-size: ${props => props.theme.fontSize};
    font-weight: bold;
  }
`

export const ErrorMessage = styled.p`
  font-size: ${props => props.theme.fontSize};
  margin: 0 10px 10px 10px;
  color: ${props => props.theme.errorColor};
  font-weight: bold;
`

export const LabelText = styled.p`
  transition: ease all 200ms;
  position: absolute;
  font-size: 20px;
  top: 10px;
  left: 10px;
  bottom: 10px;
  font-style: italic;
`

export const Label = styled.label`
  position: relative;
  width: 100%;
`

export const InputWrapper = styled.div`
  padding-top: 30px;
`

const FormikInput: FC<props> = ({placeholder, wrapperClassName, className, label, ...props}) => {
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
        <InputField
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

export default FormikInput