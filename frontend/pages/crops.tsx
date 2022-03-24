import { useState } from "react";
import type { NextPage } from "next";
import Link from "next/link";
import { Card, Col, Form, Row, Tab, Table, Tabs } from "react-bootstrap";
import { FaEdit, FaPlus } from "react-icons/fa";

import ButtonIcon from "../components/ButtonIcon";
import ModalContainer from "../components/ModalContainer";
import Layout from "../components/Layout";
import { cropsData } from "../data";
import useModal from "../hooks/useModal";

const Crops: NextPage = () => {
  const { modalOpen, showModal, closeModal } = useModal();
  const [area, setArea] = useState("");
  const [plantType, setPlantType] = useState("");
  const [variety, setVariety] = useState("");
  const [quantity, setQuantity] = useState("");
  const [type, setType] = useState("");
  const [isError, setIsError] = useState(false);

  const addCrops = () => {
    if (!area || !plantType || !variety || !quantity || !type) {
      setIsError(true);
    } else {
      setIsError(false);
      closeModal();
    }
  };

  return (
    <Layout>
      <Row>
        <Col>
          <h3 className="pb-3">Crops</h3>
        </Col>
      </Row>
      <Row>
        <Col>
          <Tabs defaultActiveKey="batch">
            <Tab eventKey="batch" title="Batch" className="pt-3">
              <ButtonIcon
                label="Add New Batch"
                icon={<FaPlus className="me-2" />}
                onClick={showModal}
                variant="primary"
              />
              <Table responsive className="my-4">
                <thead>
                  <tr>
                    <th>Crop Variety</th>
                    <th>Batch ID</th>
                    <th>Days Since Seeding</th>
                    <th>Initial Quantity</th>
                    <th>Status</th>
                    <th />
                  </tr>
                </thead>
                <tbody>
                  {cropsData &&
                    cropsData.map(
                      ({
                        id,
                        varieties,
                        batchId,
                        daysSinceSeeding,
                        qty,
                        seeding,
                        growing,
                        dumped,
                      }) => (
                        <tr key={id}>
                          <td>
                            <Link href={`/crops/${id}`}>{varieties}</Link>
                          </td>
                          <td>{batchId}</td>
                          <td>{daysSinceSeeding}</td>
                          <td>{qty}</td>
                          <td>{`${seeding} Seeding, ${growing} Growing, ${dumped} Dumped`}</td>
                          <td>
                            <FaEdit
                              onClick={showModal}
                              className="show-pointer"
                            />
                          </td>
                        </tr>
                      )
                    )}
                </tbody>
              </Table>
            </Tab>
            <Tab eventKey="archives" title="Archives">
              <Table responsive className="my-4">
                <thead>
                  <tr>
                    <th>Crop Variety</th>
                    <th>Batch ID</th>
                    <th>Days Since Seeding</th>
                    <th>Initial Quantity</th>
                    <th>Status</th>
                    <th />
                  </tr>
                </thead>
                <tbody>
                  <tr>
                    <td colSpan={6}>No crops available.</td>
                  </tr>
                </tbody>
              </Table>
            </Tab>
          </Tabs>
        </Col>
      </Row>
      <ModalContainer
        title="Update Crop"
        isShow={modalOpen}
        handleCloseModal={closeModal}
        handleSubmitModal={addCrops}
      >
        <>
          <small className="text-muted">
            Crop Batch is a quantity or consignment of crops done at one time.
          </small>
          <Form className="mt-3">
            <Form.Group className="mb-3">
              <Form.Label>Select activity type of this crop batch</Form.Label>
              <Form.Check
                name="activity"
                type="radio"
                label="Seeding"
                defaultChecked
              />
              <Form.Check name="activity" type="radio" label="Growing" />
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>Area</Form.Label>
              <Form.Select onChange={(e) => setArea(e.target.value)}>
                <option>Select Area to Grow</option>
                <option value="1">Organic Lettuce</option>
                <option value="2">Organic Chilli</option>
              </Form.Select>
              {isError && (
                <Form.Text className="text-danger">
                  The initial area field is required
                </Form.Text>
              )}
            </Form.Group>
            <Form.Group className="mb-3">
              <Row>
                <Col>
                  <Form.Label>Plant Type</Form.Label>
                  <Form.Select onChange={(e) => setPlantType(e.target.value)}>
                    <option>Select Plant Type</option>
                    <option value="1">Vegetable</option>
                  </Form.Select>
                  {isError && (
                    <Form.Text className="text-danger">
                      The plant type field is required
                    </Form.Text>
                  )}
                </Col>
                <Col>
                  <Form.Label>Crop Variety</Form.Label>
                  <Form.Select onChange={(e) => setVariety(e.target.value)}>
                    <option>Select Crop Variety</option>
                    <option value="1">Romaine</option>
                  </Form.Select>
                  {isError && (
                    <Form.Text className="text-danger">
                      The crop variety field is required
                    </Form.Text>
                  )}
                </Col>
              </Row>
            </Form.Group>
            <Form.Group className="mb-3">
              <Row>
                <Col>
                  <Form.Label>Container Quantity</Form.Label>
                  <Form.Control
                    type="text"
                    onChange={(e) => setQuantity(e.target.value)}
                  />
                  {isError && (
                    <Form.Text className="text-danger">
                      The container quantity field is required
                    </Form.Text>
                  )}
                </Col>
                <Col>
                  <Form.Label>Container Type</Form.Label>
                  <Form.Select onChange={(e) => setType(e.target.value)}>
                    <option>Select Unit</option>
                    <option value="1">Pots</option>
                    <option value="2">Trays</option>
                  </Form.Select>
                  {isError && (
                    <Form.Text className="text-danger">
                      The container type field is required
                    </Form.Text>
                  )}
                </Col>
              </Row>
            </Form.Group>
          </Form>
        </>
      </ModalContainer>
    </Layout>
  );
};

export default Crops;
