import { storiesOf } from '@storybook/vue';
import BtnBack from './btn-back.vue';

storiesOf('BtnBack', module)
  .add('default', () => (
    {
      components: { BtnBack },
      template: '<BtnBack :route="{name: \'IntroFarmCreate\'}" />',
    }
  ))
  .add('with custom css', () => (
    {
      components: { BtnBack },
      template: '<BtnBack :route="{name: \'IntroFarmCreate\'}" customClass="float-right" />',
    }
  ));
