import React, {FC, useMemo} from "react";
import {useFormikContext} from "formik";
import styled from "styled-components";

const Err = styled.p`
  font-size: ${props => props.theme.fontSize};
  color: ${props => props.theme.errorColor};
  font-weight: bold;
`

const FormikStandardError: FC = () => {
  const data = useFormikContext()
  return useMemo(() => {
    if (data?.status?.DEFAULT_ERROR) {
      return <Err>{data.status.DEFAULT_ERROR}</Err>
    }

    return null
  }, [data])
}

export default FormikStandardError