#ifndef PLURALITY_COMPONENTS_HELLOCOMPONENT_HPP
#define PLURALITY_COMPONENTS_HELLOCOMPONENT_HPP

#include <string>
#include <iostream>

#include "Component.hpp"

class HelloComponent : public Component {
	public:
		HelloComponent();
		virtual std::map<std::string, std::string> getPossibleValues() const override;
		inline virtual void Start() override;
};

HelloComponent::HelloComponent()
	: Component("HelloComponent")
{
}

std::map<std::string, std::string> HelloComponent::getPossibleValues() const
{
	return {
		{"greetee", "string"},
		{"number of greets", "int"}
	};
}

void HelloComponent::Start()
{
	const std::string& g = boost::get<std::string>(mValues["greetee"]);
	int num = boost::get<int>(mValues["number of greets"]);
	for(int i = 0; i < num; i++)
		std::cout << "Hello " << g << "!\n";
}

#endif

