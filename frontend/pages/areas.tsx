import { useState } from "react";
import type { NextPage } from "next";
import { Col, Form, Row } from "react-bootstrap";
import { FaPlus } from "react-icons/fa";

import ButtonIcon from "../components/ButtonIcon";
import ModalContainer from "../components/ModalContainer";
import PanelArea from "../components/PanelArea";
import Layout from "../components/Layout";
import { areaData } from "../data";
import useModal from "../hooks/useModal";

const Areas: NextPage = () => {
  const { modalOpen, showModal, closeModal } = useModal();
  const [areaName, setAreaName] = useState("");
  const [areaSize, setAreaSize] = useState("");
  const [areaSizeUnit, setAreaSizeUnit] = useState("Ha");
  const [areaType, setAreaType] = useState("Seeding");
  const [areaLocations, setAreaLocations] = useState("Field (Outdoor");
  const [reservoir, setReservoir] = useState("");
  const [isError, setIsError] = useState(false);

  const addAreas = () => {
    if (!areaName || !areaSize || !reservoir) {
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
          <h3 className="pb-3">Areas</h3>
        </Col>
      </Row>
      <Row>
        <Col>
          <ButtonIcon
            label="Add Area"
            icon={<FaPlus className="me-2" />}
            onClick={showModal}
            variant="primary"
          />
        </Col>
      </Row>
      <Row className="mt-3">
        {areaData.map(
          ({ id, name, type, size, unit, batches, quantity, edit }) => (
            <Col xs={12} sm={12} md={6} lg={4} key={id}>
              <PanelArea
                id={id}
                name={name}
                type={type}
                size={size}
                unit={unit}
                batches={batches}
                quantity={quantity}
                edit={edit}
                onClick={showModal}
              />
            </Col>
          )
        )}
      </Row>
      <ModalContainer
        title="Add New Area"
        isShow={modalOpen}
        handleCloseModal={closeModal}
        handleSubmitModal={addAreas}
      >
        <>
          <small className="text-muted">
            Area is a space where you grow your plants. It could be a seeding
            tray, a garden bed, or a pot or anything that describes the
            different physical locations in your facility.
          </small>
          <Form className="mt-3">
            <Form.Group className="mb-3">
              <Form.Label>Area Name</Form.Label>
              <Form.Control
                type="text"
                onChange={(e) => setAreaName(e.target.value)}
              />
              {isError && (
                <Form.Text className="text-danger">
                  The name field is required
                </Form.Text>
              )}
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>Size</Form.Label>
              <Row>
                <Col>
                  <Form.Control
                    type="text"
                    onChange={(e) => setAreaSize(e.target.value)}
                  />
                  {isError && (
                    <Form.Text className="text-danger">
                      The size field is required
                    </Form.Text>
                  )}
                </Col>
                <Col>
                  <Form.Select
                    onChange={(e) => setAreaSizeUnit(e.target.value)}
                  >
                    <option value="Ha">Ha</option>
                    <option value="m2">m2</option>
                  </Form.Select>
                </Col>
              </Row>
            </Form.Group>
            <Form.Group className="mb-3">
              <Row>
                <Col xs={4}>
                  <Form.Label>Type</Form.Label>
                  <Form.Select onChange={(e) => setAreaType(e.target.value)}>
                    <option value="Seeding">Seeding</option>
                    <option value="Growing">Growing</option>
                  </Form.Select>
                </Col>
                <Col xs={8}>
                  <Form.Label>Locations</Form.Label>
                  <Form.Select
                    onChange={(e) => setAreaLocations(e.target.value)}
                  >
                    <option value="Field (Outdoor)">
                      Field Outdoor (Outdoor)
                    </option>
                    <option value="Greenhouse (Indoor)">
                      Greendhouse (Indoor)
                    </option>
                  </Form.Select>
                </Col>
              </Row>
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>Select Reservoir</Form.Label>
              <Form.Select onChange={(e) => setReservoir(e.target.value)}>
                <option value="">Please select reservoir</option>
                <option value="River">River</option>
                <option value="Water Tank">Water Tank</option>
              </Form.Select>
              {isError && (
                <Form.Text className="text-danger">
                  The reservoir field is required
                </Form.Text>
              )}
            </Form.Group>
            <Form.Group>
              <Form.Label>
                Select photo <small className="text-muted">(optional)</small>
              </Form.Label>
              <Form.Control type="file" />
            </Form.Group>
          </Form>
        </>
      </ModalContainer>
    </Layout>
  );
};

export default Areas;
