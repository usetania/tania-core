import { useState } from "react";
import type { NextPage } from "next";
import Image from "next/image";
import { useRouter } from "next/router";
import { Card, Col, Container, Form, Row } from "react-bootstrap";
import { FaUnlock } from "react-icons/fa";

import ButtonIcon from "@/components/ButtonIcon";

const SignIn: NextPage = () => {
  const router = useRouter();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [usernameIsEmpty, setUsernameIsEmpty] = useState(false);
  const [passwordIsEmpty, setPasswordIsEmpty] = useState(false);

  const handleSubmit = () => {
    if (!username || !password) {
      setUsernameIsEmpty(!username);
      setPasswordIsEmpty(!password);
    } else {
      router.push("/");
    }
  };

  return (
    <div className="bg-gray d-flex align-items-center vh-100">
      <Container fluid>
        <Row>
          <Col sm={12} md={{ span: 6, offset: 3 }} lg={{ span: 4, offset: 4 }}>
            <Card>
              <Card.Body>
                <div className="text-center">
                  <Image alt="Tania logo" src={"/img/logo.png"} width={200} height={65} />
                </div>
                <Form>
                  <Form.Group className="mb-3">
                    <Form.Label>Username</Form.Label>
                    <Form.Control
                      type="text"
                      placeholder="Input your username here"
                      onChange={(e) => setUsername(e.target.value)}
                    />
                    {usernameIsEmpty && (
                      <Form.Text className="text-danger">
                        The username field is required
                      </Form.Text>
                    )}
                  </Form.Group>
                  <Form.Group className="mb-3">
                    <Form.Label>Password</Form.Label>
                    <Form.Control
                      type="password"
                      placeholder="Your password here"
                      onChange={(e) => setPassword(e.target.value)}
                    />
                    {passwordIsEmpty && (
                      <Form.Text className="text-danger">
                        The password field is required
                      </Form.Text>
                    )}
                  </Form.Group>
                </Form>
                <div className="text-center">
                  <ButtonIcon
                    label="Login"
                    icon={<FaUnlock className="me-1" />}
                    type="submit"
                    variant="primary"
                    onClick={handleSubmit}
                  />
                </div>
              </Card.Body>
            </Card>
          </Col>
        </Row>
      </Container>
    </div>
  );
}

export default SignIn;