import Link from "next/link";
import { Card, Col, Row } from "react-bootstrap";
import { FaEdit } from "react-icons/fa";

interface iPanelArea {
  id: number;
  name: string;
  type: string;
  size: number;
  unit: string;
  batches: number;
  quantity: number;
  edit: boolean;
  onClick: () => void;
}

const PanelArea = ({
  id,
  name,
  type,
  size,
  unit,
  batches,
  quantity,
  edit,
  onClick,
}: iPanelArea): JSX.Element => {
  return (
    <Card className="mb-3">
      <Card.Body>
        <Card.Title>
          <Row>
            <Col xs={8}>
              <h5>
                <Link href={`/areas/${id}`}>{name}</Link>
              </h5>
            </Col>
            <Col xs={4}>
              <div className="d-flex justify-content-end">
                {edit && (
                  <FaEdit className="me-2 show-pointer" onClick={onClick} />
                )}
              </div>
            </Col>
          </Row>
        </Card.Title>
        <small className="text-muted">{type}</small>
        <Row className="mt-3">
          <Col xs={4}>
            <small className="text-muted">{`Size ${unit}`}</small>
            <br />
            {size}
          </Col>
          <Col xs={4}>
            <small className="text-muted">Batches</small>
            <br />
            {batches}
          </Col>
          <Col xs={4}>
            <small className="text-muted">Quantity</small>
            <br />
            {quantity}
          </Col>
        </Row>
      </Card.Body>
    </Card>
  );
};

export default PanelArea;
