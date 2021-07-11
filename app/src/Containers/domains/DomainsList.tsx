import React, {FC} from "react";
import {ContentBox, PageWrapper, SectionHeader} from "../../Styled";
import useDomains, {domainI} from "../../Hooks/DomainsHook";
import ContentLoader from "../../Components/ContentLoader";
import styled from "styled-components";
import Button from "../../Components/Button";
import {useHistory} from "react-router-dom";

const DomainsList: FC = () => {
  const history = useHistory()
  const {query} = useDomains()

  return (
    <PageWrapper>
      <ContentBox>
        <SectionHeader>Domains</SectionHeader>
        <Button
          onClick={() => history.push("/domains/new")}
          style={{marginBottom: 20}}
        >
          New Domain
        </Button>
        <ContentLoader info={query}>
          {query.data && query.data.map(el => {
            return <Line key={`domain-${el.id}`} data={el}/>
          })}
        </ContentLoader>
      </ContentBox>
    </PageWrapper>
  )
}

const LineDiv = styled.div`
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: ${props => props.theme.standardShadow};
  border-radius: ${props => props.theme.standardBorderRadius};
  padding: 10px 20px;
  margin-bottom: 25px;

  .green {
    color: ${props => props.theme.successColor};
  }

  .orange {
    color: ${props => props.theme.orangeColor};
  }
`

const Status = styled.span`
  display: inline-block;
  width: 100px;
`

const Line: FC<{ data: domainI }> = ({data}) => {
  const history = useHistory()

  return (
    <LineDiv>
      <div>
        <Status
          className={data.verified === 0 ? "orange" : "green"}>{data.verified === 0 ? "Not Verified" : "Verified"}</Status>
        {data.domain_url}
      </div>
      <div>
        <Button onClick={() => history.push(`/domains/${data.id}`)}>View</Button>
      </div>
    </LineDiv>
  )
}

export default DomainsList;