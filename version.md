# Version history

### 1.1
Metric names changed;
* charged_energy renamed to charged_energy_Wh
* charing_state to charching_state
  
Calculation issue fixed: 
* charged_energy_Wh represents Modbus TCP 1502. This returns the wrong factor value for Wh. For this reason divided by 10.
* total_energy_counter_Wh represents Modbus TCP 1036. This returns the wrong factor value for Wh. For this reason divided by 10.


### 1.0.1
Documentation fixed.

### 1.0
Initial version
