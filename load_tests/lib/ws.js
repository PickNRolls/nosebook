import { sleep } from 'k6';

const wsRampingVus = (stages) => {
  let vus = 0;

  for (let i = 0; i < stages.length; i++) {
    if (stages[i].target > vus) {
      vus = stages[i].target;
    }
  }

  return vus;
};

const calcVuDuration = (vuIndex, stages) => {
  let duration = 0;
  const vus = wsRampingVus(stages);

  for (let i = 1; i < stages.length - 1; i++) {
    duration += stages[i].duration;
  }

  const rampUpStage = stages[0];
  duration += rampUpStage.duration / vus * (vus - vuIndex)

  const rampDownStage = stages[stages.length - 1];
  duration += rampDownStage.duration / vus * (vuIndex + 1);

  return duration;
};

const calcTotalDuration = (stages) => {
  let duration = 0;

  for (let i = 0; i < stages.length; i++) {
    duration += stages[i].duration;
  }

  return duration;
}

export const rampUpWsOptions = (stages) => {
  let vus = wsRampingVus(stages);
  const rampUpDuration = stages[0].duration;
  const duration = calcTotalDuration(stages);

  return {
    scenarios: {
      wsRamping: {
        executor: 'per-vu-iterations',
        vus,
        iterations: 1,
        maxDuration: `${(duration + rampUpDuration) / 1000}s`,
      },
    },
  };
};

const wsRampingUpSleepTimeSeconds = ({ vuIndex, vus, rampUpDuration }) => {
  return rampUpDuration / vus * vuIndex / 1000;
};

export const rampUpWs = (stages) => {
  const vuIndex = __VU - 1;
  const vus = wsRampingVus(stages);
  const rampUpDuration = stages[0].duration;
  const duration = calcVuDuration(vuIndex, stages);
  sleep(wsRampingUpSleepTimeSeconds({
    vus,
    vuIndex,
    rampUpDuration,
  }));

  return {
    vuIndex,
    vus,
    duration,
  };
};

