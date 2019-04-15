import { storiesOf } from '@storybook/vue';
import BtnWatering from './btn-watering.vue';

storiesOf('BtnWatering', module)
  .add('default', () => (
    {
      components: { BtnWatering },
      template: '<BtnWatering />',
    }
  ))
  .add('with custom css', () => (
    {
      components: { BtnWatering },
      template: '<BtnWatering class="float-right" />',
    }
  ));
