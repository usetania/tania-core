import { useState } from "react";
import type { NextPage } from "next";
import Image from "next/image";
import { useRouter } from "next/router";
import { Card, Col, Container, Form, Row } from "react-bootstrap";
import { FaUserPlus } from "react-icons/fa";

import ButtonIcon from "@/components/ButtonIcon";

const Register: NextPage = () => {
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
                  <h4 className="py-3">Let&apos;s setup your account</h4>
                </div>
                <Form>
                  <Form.Group className="mb-3">
                    <Form.Label>Username</Form.Label>
                    <Form.Control
                      type="text"
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
                      onChange={(e) => setPassword(e.target.value)}
                    />
                    {passwordIsEmpty && (
                      <Form.Text className="text-danger">
                        The password field is required
                      </Form.Text>
                    )}
                  </Form.Group>
                </Form>
                <div className="text-center mt-4">
                  <ButtonIcon
                    label="Create admin account"
                    icon={<FaUserPlus className="me-1" />}
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
};

export default Register;
