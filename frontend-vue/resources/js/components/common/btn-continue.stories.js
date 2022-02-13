import { storiesOf } from '@storybook/vue';
import BtnContinue from './btn-continue.vue';

storiesOf('BtnContinue', module)
  .add('default', () => (
    {
      components: { BtnContinue },
      template: '<BtnContinue title="Continue" />',
    }
  ))
  .add('default finish', () => (
    {
      components: { BtnContinue },
      template: '<BtnContinue title="Finish Setup" />',
    }
  ))
  .add('with custom css', () => (
    {
      components: { BtnContinue },
      template: '<BtnContinue title="Continue" class="float-right" />',
    }
  ));
