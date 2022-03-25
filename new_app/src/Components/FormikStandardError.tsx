import React, {FC, useMemo} from "react";
import {useFormikContext} from "formik";
import styled from "styled-components";

const Err = styled.p<{spaceTop?: boolean}>`
  font-size: ${props => props.theme.fontSize};
  color: ${props => props.theme.errorColor};
  font-weight: bold;
  
  ${props => props.spaceTop && `margin-top: 20px;`}
`

const FormikStandardError: FC<{spaceTop?: boolean}> = props => {
  const data = useFormikContext()
  return useMemo(() => {
    if (data?.status?.DEFAULT_ERROR) {
      return <Err spaceTop={props.spaceTop}>{data.status.DEFAULT_ERROR}</Err>
    }

    return null
  }, [data])
}

export default FormikStandardError