import React, {FC, useMemo, useState} from "react";
import Select, {ValueType} from 'react-select';
import _ from "lodash";
import {ErrorMessage, InputWrapper, Label, LabelText} from "./FormikInput";
import {useField, useFormikContext} from "formik";
import styled, {useTheme} from "styled-components";
import {themeType} from "../Styled";

interface propsI {
  wrapperClassName?: string
  name: string
  label: string
  className?: string
  isMulti?: boolean
  options: Option[]
}

interface Option {
  label: string
  value: string
}

const DownArrow = styled.div`
  width: 0;
  height: 0;
  border-left: 6px solid transparent;
  border-right: 6px solid transparent;
  margin-right: 10px;

  border-top: 6px solid ${props => props.theme.black};
`

const SelectWrapper = styled(InputWrapper)`
  font-size: 16px;
  
  .f-select__menu {
    margin-top: 0;
    box-shadow: none;
  }

  .f-select__indicator-separator {
    display: none;
  }

  .f-select__control {
    height: 42px;
    border: 3px solid ${props => props.theme.black};
    border-radius: 0;
    border-top: 0;
    border-right: 0;
    border-left: 0;
    transition: ease all 200ms;
    margin-bottom: 20px;

    &:hover {
      border-bottom: 3px solid ${props => props.theme.black} !important;
    }

    &--is-focused {
      border: 3px solid ${props => props.theme.black} !important;
      border-radius: 0 !important;
      border-top: 0 !important;
      border-right: 0 !important;
      border-left: 0 !important;
      box-shadow: none;
    }
  }

  .formik-sel {
    &.is-focused + .label-holder {
      top: 10px;
      font-size: ${props => props.theme.fontSize};
      font-weight: bold;
    }

    &.has-val + .label-holder {
      top: 10px;
      font-size: ${props => props.theme.fontSize};
      font-weight: bold;
    }
  }

  .f-select__option {
    &:hover {
      cursor: pointer;
    }

    &--is-selected {
      background: ${props => props.theme.black};
    }
  }
`

const FormikSelect: FC<propsI> = (props) => {
  const styledTheme = useTheme() as themeType
  const {
    label,
    wrapperClassName,
    isMulti = false,
    className,
    options
  } = props

  const {status} = useFormikContext()
  const [field, meta, helpers] = useField(props)
  const [isSelected, setIsSelected] = useState(false)

  const errMsg = useMemo(() => {
    if (status && status[props.name]) {
      return <ErrorMessage>status[props.name]</ErrorMessage>
    }

    return meta.touched && meta.error && <ErrorMessage>{meta.error}</ErrorMessage>
  }, [meta, status])

  const onChange = (option: ValueType<Option | Option[]>) => {
    if (!option) {
      helpers.setValue("")
      return
    }
    helpers.setValue(
      isMulti
        ? (option as Option[]).map((item: Option) => item.value)
        : (option as Option).value
    );
  };

  const getValue = () => {
    if (options) {
      if (isMulti) {
        const r = options.filter(option => field.value.indexOf(option.value) >= 0)
        if (r) {
          return r
        }
      } else {
        const r = options.find(option => option.value === field.value);
        if (r) {
          return r
        }
      }

      return null
    } else {
      return isMulti ? [] : ("" as any);
    }
  };

  const handleBlur = (e: React.FocusEvent<any>) => {
    setIsSelected(false)
    field.onBlur(e)
  }

  return (
    <SelectWrapper className={wrapperClassName ? wrapperClassName : ""}>
      <Label>
        {/*<InputField*/}
        {/*  className={`${className} ${errMsg && "error"} ${meta.touched && "touched"} ${!_.isEmpty(field.value) && "has-val"}`}*/}
        {/*/>*/}
        <Select
          className={`${className} ${errMsg && "error"} ${meta.touched && "touched"} ${!_.isEmpty(field.value) && "has-val"} ${isSelected && "is-focused"} formik-sel`}
          classNamePrefix={"f-select"}
          onBlur={handleBlur}
          onFocus={() => setIsSelected(true)}
          placeholder={""}
          options={options}
          value={getValue()}
          onChange={onChange}
          theme={theme => {
            return {
              ...theme,
              borderRadius: 0,
              colors: {
                ...theme.colors,
                primary25: styledTheme.lightBlack,
                primary50: styledTheme.lightBlack,
                primary: styledTheme.lightBlack,
              },
            }
          }}
          components={{
            DropdownIndicator: () => {
              return <DownArrow />
            }
          }}
        />
        {label && <LabelText className={"label-holder"}>{label}</LabelText>}
        {errMsg}
      </Label>
    </SelectWrapper>
  )
}

export default FormikSelect;