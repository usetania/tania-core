import { useState } from "react";
import type { NextPage } from "next";
import { Card, Col, Form, Row } from "react-bootstrap";
import { FaCheck } from "react-icons/fa";

import ButtonIcon from "../components/ButtonIcon";
import Layout from "../components/Layout";

const Account: NextPage = () => {
  const [username, setUsername] = useState("tania");
  const [oldPassword, setOldPassword] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [isError, setIsError] = useState(false);

  const changePassword = () => {
    if (!oldPassword || !password) {
      setIsError(true);
    } else {
      setIsError(false);
    }
  };

  return (
    <Layout>
      <Row>
        <Col>
          <h3 className="pb-3">Account Settings</h3>
        </Col>
      </Row>
      <Row>
        <Col xs={12} sm={12} md={6}>
          <Card>
            <Card.Body>
              <Form>
                <Form.Group className="mb-3">
                  <Form.Label>Username</Form.Label>
                  <Form.Control type="text" readOnly value={username} />
                </Form.Group>
                <Form.Group className="mb-3">
                  <Form.Label>Old Password</Form.Label>
                  <Form.Control
                    type="password"
                    value={oldPassword}
                    onChange={(e) => setOldPassword(e.target.value)}
                  />
                  {isError && (
                    <Form.Text className="text-danger">
                      The old password field is required
                    </Form.Text>
                  )}
                </Form.Group>
                <Form.Group className="mb-3">
                  <Form.Label>Password</Form.Label>
                  <Form.Control
                    type="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                  />
                  {isError && (
                    <Form.Text className="text-danger">
                      The password field is required
                    </Form.Text>
                  )}
                </Form.Group>
                <Form.Group className="mb-3">
                  <Form.Label>Confirm Password</Form.Label>
                  <Form.Control
                    type="password"
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                  />
                </Form.Group>
                <ButtonIcon
                  label="Save"
                  icon={<FaCheck className="me-1" />}
                  variant="secondary"
                  onClick={changePassword}
                  textColor="text-light"
                />
              </Form>
            </Card.Body>
          </Card>
        </Col>
      </Row>
    </Layout>
  );
};

export default Account;
