/**
 * Frontend Payment Methods Testing Script
 * Run this in the browser console on the payment settings page
 * to test the payment methods functionality
 */

async function testPaymentMethodsFrontend() {
    console.log('🚀 Starting Payment Methods Frontend Testing...');
    
    try {
        // Test 1: Check if payment methods are loaded
        console.log('📋 Test 1: Checking if payment methods are loaded...');
        
        // Wait for payment methods to load
        await new Promise(resolve => setTimeout(resolve, 2000));
        
        // Look for payment method cards
        const paymentMethodCards = document.querySelectorAll('[class*="border rounded-md"]');
        console.log(`✅ Found ${paymentMethodCards.length} payment method cards`);
        
        // Test 2: Check for toggle checkboxes
        console.log('🔄 Test 2: Checking toggle checkboxes...');
        const toggleCheckboxes = document.querySelectorAll('input[type="checkbox"][id*="enable-"]');
        console.log(`✅ Found ${toggleCheckboxes.length} toggle checkboxes`);
        
        // Test 3: Check for form inputs
        console.log('📝 Test 3: Checking form inputs...');
        const formInputs = document.querySelectorAll('input[type="text"], input[type="password"]');
        console.log(`✅ Found ${formInputs.length} form inputs`);
        
        // Test 4: Check for save buttons
        console.log('💾 Test 4: Checking save buttons...');
        const saveButtons = document.querySelectorAll('button[type="submit"]');
        console.log(`✅ Found ${saveButtons.length} save buttons`);
        
        // Test 5: Display current payment method status
        console.log('📊 Test 5: Current payment method status:');
        paymentMethodCards.forEach((card, index) => {
            const titleElement = card.querySelector('h3');
            const checkboxElement = card.querySelector('input[type="checkbox"][id*="enable-"]');
            
            if (titleElement && checkboxElement) {
                const title = titleElement.textContent;
                const isEnabled = checkboxElement.checked;
                console.log(`   ${title}: ${isEnabled ? '✅ Enabled' : '❌ Disabled'}`);
            }
        });
        
        console.log('🎉 Frontend testing completed successfully!');
        return true;
        
    } catch (error) {
        console.error('❌ Frontend testing failed:', error);
        return false;
    }
}

// Auto-run the test
testPaymentMethodsFrontend();
