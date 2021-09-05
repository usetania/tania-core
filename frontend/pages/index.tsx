import type { NextPage } from "next";
import { Button, Col, Container, Row } from "react-bootstrap";

const Home: NextPage = () => {
  return (
    <>
      <main>
        <Container>
          <Row>
            <Col>
              <Button variant="primary">Primary</Button>
            </Col>
          </Row>
        </Container>
      </main>
    </>
  );
};

export default Home;
