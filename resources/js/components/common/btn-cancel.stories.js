import { storiesOf } from '@storybook/vue';
import BtnCancel from './btn-cancel.vue';

storiesOf('BtnCancel', module)
  .add('default', () => (
    {
      components: { BtnCancel },
      template: '<BtnCancel />',
    }
  ))
  .add('with custom css', () => (
    {
      components: { BtnCancel },
      template: '<BtnCancel customClass="float-right" />',
    }
  ));
