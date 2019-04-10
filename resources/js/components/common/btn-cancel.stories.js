import { storiesOf } from '@storybook/vue';
import BtnCancel from './btn-cancel.vue';

storiesOf('BtnCancel', module)
  .add('default', () => (
    {
      components: { BtnCancel },
      template: '<BtnCancel />',
    }
  ))
  .add('float-right', () => (
    {
      components: { BtnCancel },
      template: '<BtnCancel customClass="float-right" />',
    }
  ));
