
type TPackage = {
    type: string;
    data: TFinding | TObservation;
};

type TFinding = {
    name: string;
    patient: string;
    sensor: string;
    value: string;
    timestamp: string;
};

type TObservation = {
    name: string;
    sensor: string;
    value: string;
    timestamp: string;
};

export type { TPackage, TFinding, TObservation };